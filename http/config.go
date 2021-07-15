package http

import "github.com/spf13/viper"

var config = viper.New()

const serverShutdownTimeoutConfig = "server_shutdown_timeous"

func init() {
	config.AutomaticEnv()
	config.SetDefault(serverShutdownTimeoutConfig, "10s")
}
