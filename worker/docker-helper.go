package main

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func runDockerContainer() io.ReadCloser {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	bash := `
	#!/bin/bash
	ls -la
	hostname
	echo "sleeping ..."
	sleep 2s
	pwd
	sleep 2s
	date
	`

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		// Cmd:   []string{"/bin/sh"},
		Cmd:          []string{"/bin/sh", "-c", "touch hello.sh && echo " + bash + " > hello.sh && chmod +x hello.sh && ./hello.sh"},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   []string{},
	}, nil, nil, "")

	cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		// ShowStderr: true,
		Follow: true,
		// Details:    true,
		// Tail:       "100",
	})
	if err != nil {
		panic(err)
	}

	return out

}
