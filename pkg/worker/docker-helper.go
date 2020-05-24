package worker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
	"workhorse/api"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func createTempFile(contents string, jobFile string) string {
	jobDir, err := ioutil.TempDir("", "job")
	if err != nil {
		fmt.Println("Directory error", err)
		return ""
	}

	fullPath := path.Join(jobDir, jobFile)
	fmt.Println("Creating file::" + fullPath)

	f, err := os.Create(fullPath)
	os.Chmod(fullPath, 0777)
	if err != nil {
		fmt.Println("File creation error", err)
		return ""
	}

	f.Write([]byte(contents))
	f.Close()

	return jobDir
}

func getTargetMountDirectory() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func runDockerContainer(job *api.JobTransferObject) io.ReadCloser {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// reader, err := cli.ImagePull(ctx, "docker.io/library/"+job.Image, types.ImagePullOptions{})
	reader, err := cli.ImagePull(ctx, job.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	jobDir := createTempFile(string(job.ScriptContents), job.FileName)
	fmt.Println(jobDir)

	targetMountDir := getTargetMountDirectory()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        job.Image,
		Cmd:          []string{"/bin/sh", "-c", "./" + targetMountDir + "/" + job.FileName},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   []string{},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: jobDir,
				//Use a unique different directory to avoid conflict with directory name
				//in container
				// Target: "/job",
				Target: path.Join("/", targetMountDir),
			},
		},
	}, nil, "")

	cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	return out
}
