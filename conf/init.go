package conf

import (
	"context"
	"errors"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	PushGatewayEndpoint string `env:"PUSHGATEWAY_ENDPOINT"`
	BindAddr            string `env:"BIND_ADDR,default=0.0.0.0:19091"`
}

var cfg config

func init() {
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		panic(err)
	}

	if cfg.PushGatewayEndpoint == "" {
		panic(errors.New("the environment variable PUSH_GATEWAY_ENDPOINT is required"))
	}
}

func PushGatewayEndpoint() string {
	return cfg.PushGatewayEndpoint
}

func BindAddr() string {
	return cfg.BindAddr
}
