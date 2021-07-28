package app

import (
	"github.com/7d4b9/lever/example/tests/docker"
	"github.com/docker/docker/api/types/container"
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

func NewClient(client *docker.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) NewContainer() error {
	return c.client.RunContainer(&container.Config{
		Env: []string{
			"HTTP_PORT=8088",
		},
	}, config.GetString(imageConfig))
}
