package app

import (
	"fmt"

	"ricklepickle/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

var config = viper.New()

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("example")
}

const imageConfig = "image"

type Client struct {
	client *docker.Client
}

func NewClient(client *client.Client) (*Client, error) {
	image := config.GetString(imageConfig)
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
			"HTTP_PORT=8088",
		},
	}, config.GetString(imageConfig))
}
