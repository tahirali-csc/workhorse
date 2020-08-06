package worker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
	"workhorse/pkg/api"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func createTempFile(contents string, jobFile string) string {
	jobDir, err := ioutil.TempDir("", "job")
	if err != nil {
		log.Println("Directory error", err)
		return ""
	}

	jobDir, err = filepath.EvalSymlinks(jobDir)
	if err != nil {
		log.Println("Symlink error", err)
		return ""
	}

	fullPath := path.Join(jobDir, jobFile)
	log.Println("Creating file::" + fullPath)

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

func runDockerContainer(job *api.WorkflowJob) io.ReadCloser {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	log.Printf("Pulling an image %s", job.Image)
	reader, err := cli.ImagePull(ctx, job.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	jobDir := createTempFile(string(job.ScriptContents), job.FileName)
	log.Println(jobDir)

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

	if err != nil {
		log.Fatal(err)
		return nil
	}

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
