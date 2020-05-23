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

	jobDir := createTempFile(string(job.ScriptContents), "dumb.sh")
	fmt.Println(jobDir)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"/bin/sh", "-c", "./job/dumb.sh"},
		// Cmd: []string{"ls -la"},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   []string{},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				// Source: "//home/tahir/workspace/rnd-projects/workhorse/",
				// Target: "/java",
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
