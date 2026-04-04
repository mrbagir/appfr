package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type ConfigTest struct {
	AppEnv     string `env:"APP_ENV"   envDefault:"dev"       validate:"required"`
	HttpPort   int    `env:"HTTP_PORT" envDefault:"8080"      validate:"required"`
	GrpcConfig ConfigTest2
}

type ConfigTest2 struct {
	Domain   string `env:"DOMAIN"    envDefault:"localhost" validate:"required"`
	GrpcPort int    `env:"GRPC_PORT" envDefault:"9090"      validate:"required"`
}

func TestFunction(t *testing.T) {
	config := &ConfigTest{}

	a := &App{}
	a.ParseConfig(config)
	fmt.Printf("%+v", config)
}

func TestParseEnv(t *testing.T) {
	type Inner struct {
		A string `env:"OLA" envDefault:"HI"`
	}
	type Config struct {
		NilInner  *Inner
		InitInner *Inner `env:",init"`
	}
	var cfg Config

	a := &App{}
	a.ParseConfig(&cfg)
	fmt.Print(cfg.NilInner, cfg.InitInner)
	// Output: <nil> &{HI}
}

func TestAppentError(t *testing.T) {
	var err error
	err = errors.Join(err, fmt.Errorf("first error"))
	err = errors.Join(err, fmt.Errorf("second error"))

	fmt.Println(err)
	// Output: first error; second error
}

func BenchmarkReadStoreValues(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		value := "value"
		for i := 0; i < b.N; i++ {
			_ = value
		}
	})

	b.Run("map", func(b *testing.B) {
		m := make(map[string]string)
		m["key"] = "value"
		for i := 0; i < b.N; i++ {
			_ = m["key"]
		}
	})

	b.Run("struct", func(b *testing.B) {
		type S struct{ value string }
		var s = S{value: "value"}
		for i := 0; i < b.N; i++ {
			_ = s.value
		}
	})

	b.Run("ctx", func(b *testing.B) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")
		for i := 0; i < b.N; i++ {
			_ = ctx.Value("key").(string)
		}
	})

	b.Run("pointer", func(b *testing.B) {
		value := "value"
		ptr := &value
		for i := 0; i < b.N; i++ {
			_ = *ptr
		}
	})

	b.Run("interface", func(b *testing.B) {
		var v any = "value"
		for i := 0; i < b.N; i++ {
			v, ok := v.(string)
			if ok {
				_ = v
			}
		}
	})

	b.Run("channel", func(b *testing.B) {
		ch := make(chan string, 1)
		ch <- "value"
		for i := 0; i < b.N; i++ {
			_ = <-ch
			ch <- "value"
		}
	})
}
