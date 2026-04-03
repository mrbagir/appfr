package config

import (
	"os"
	"testing"
)

func BenchmarkGetenv(b *testing.B) {
	b.Run("os", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = os.Getenv("APP_ENV")
		}
	})

	var appenv string = os.Getenv("APP_ENV")
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = appenv
		}
	})
}
