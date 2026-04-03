package app

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/mrbagir/qcash-appcore/pkg/config"
)

func (a *App) ParseEnv(s any) {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		log.Fatal("ParseEnv: argument must be a pointer to a struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := range rt.NumField() {
		field := rt.Field(i)
		fv := rv.Field(i)

		if !fv.CanSet() {
			continue
		}

		key := field.Tag.Get("env")
		if key == "" {
			continue
		}

		val, ok := os.LookupEnv(key)
		if !ok {
			val = field.Tag.Get("envDefault")
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				log.Fatalf("ParseEnv: invalid int for %s=%q: %v", key, val, err)
			}
			fv.SetInt(n)
		}
	}
}

func (a *App) readConfig() config.Config {
	return config.Config(nil)
}
