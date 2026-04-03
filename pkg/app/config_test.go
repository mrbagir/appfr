package app

import (
	"fmt"
	"testing"
)

type ConfigTest struct {
	AppEnv   string `env:"APP_ENV"   envDefault:"dev"       validate:"required"`
	Domain   string `env:"DOMAIN"    envDefault:"localhost" validate:"required"`
	HttpPort int    `env:"HTTP_PORT" envDefault:"8080"      validate:"required"`
	GrpcPort int    `env:"GRPC_PORT" envDefault:"8080"      validate:"required"`
}

func TestFunction(t *testing.T) {
	config := &ConfigTest{}

	a := &App{}
	a.ParseEnv(config)
	fmt.Printf("%v", config)
}
