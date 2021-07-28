package docker

import (
	"fmt"
	"lever/context"
	"lever/rand"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var config = viper.New()

const (
	NetworkModeConfig = "network_mode"
)

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("docker")
	config.SetDefault(NetworkModeConfig, "host")
}

type Client struct {
	*client.Client
	Image        string
	containerID  string
	runWaitGroup sync.WaitGroup
}

func NewClient(client *client.Client, image string) (*Client, error) {
	if _, err := client.ImagePull(context.CancelOnSigInterrupt, image, types.ImagePullOptions{}); err != nil {
		return nil, fmt.Errorf("pull docker image=%s, error=%w", image, err)
	}
	return &Client{
		Client: client,
		Image:  image,
	}, nil
}

var runID = rand.StringRunes(10)

func (c *Client) RunContainer(containerConfig *container.Config, containerNameBase string) error {
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(config.GetString(NetworkModeConfig)),
	}
	ctx := context.CancelOnSigInterrupt
	containerName := containerNameBase + "-" + runID
	createdContainer, err := c.Client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("cannot create container, error=%w", err)
	}
	err = c.Client.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("cannot start container, error=%w", err)
	}
	logs, err := c.Client.ContainerLogs(ctx, createdContainer.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return fmt.Errorf("cannot configure container logs, error=%w", err)
	}
	c.runWaitGroup.Add(1)
	go func() {
		defer c.runWaitGroup.Done()
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, logs); err != nil {
			zap.L().Error("container logs demux",
				zap.Error(err),
				zap.String("container", containerName))
		}
	}()
	c.containerID = createdContainer.ID
	return nil
}

func (c *Client) WaitContainer() {
	c.runWaitGroup.Wait()
}