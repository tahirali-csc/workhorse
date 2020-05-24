package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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

func runDockerContainer(job *api.JobTransferObject) io.ReadCloser {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/"+job.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	jobDir := createTempFile(string(job.ScriptContents), job.FileName)
	fmt.Println(jobDir)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		// Image: "alpine",
		Image:        job.Image,
		Cmd:          []string{"/bin/sh", "-c", "./job/" + job.FileName},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   []string{},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: jobDir,
				Target: "/job",
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
