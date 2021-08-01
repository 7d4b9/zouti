package postgres

import (
	"fmt"

	"ricklepickle/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const image = "docker.io/library/postgres:9.6.22-alpine3.14"

type Client struct {
	client *docker.Client
}

func NewClient(client *client.Client) (*Client, error) {
	c, err := docker.NewClient(client, image)
	if err != nil {
		return nil, fmt.Errorf("app new docker client for image '%s': %w", image, err)
	}
	return &Client{
		client: c,
	}, nil
}

func (c *Client) NewContainer() error {
	return c.client.RunContainer(&container.Config{
		Env: []string{
			"POSTGRES_PASSWORD=mysecretpassword",
		},
	}, "tests-postgres")
}
