package tests

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io/ioutil"
	"net"
	"testing"
)

// Container tracks information about the docker container started for tests.
type Container struct {
	ID   string
	Host string
}

func dumpContainerLogs(t *testing.T, id string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	readCloser, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(readCloser)
	if err != nil {
		panic(err)
	}

	t.Logf("Logs for %s\n%s:", id, out)
}

func startContainer(t *testing.T, image string, port string, envs []string) *Container {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{}); err != nil {
		panic(err)
	} else {
		ioutil.ReadAll(reader)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		Env:   envs,
	}, &container.HostConfig{PublishAllPorts: true}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	id := resp.ID

	containerJson, _ := cli.ContainerInspect(ctx, id)

	endpoint := containerJson.NetworkSettings.Ports[nat.Port(port+"/tcp")][0]

	c := Container{
		ID:   id,
		Host: net.JoinHostPort(endpoint.HostIP, endpoint.HostPort),
	}

	t.Logf("Image:          %s", image)
	t.Logf("DB ContainerID: %s", c.ID)
	t.Logf("Host:            %s", c.Host)
	return &c
}

// stopContainer stops and removes the specified container.
func stopContainer(t *testing.T, id string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStop(ctx, id, nil); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}
	t.Log("Stopped:", id)

	if err := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}
	t.Log("Removed:", id)
}
