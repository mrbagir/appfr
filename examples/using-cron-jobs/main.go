package main

import (
	"fmt"
	"sync"

	appcore "github.com/mrbagir/qcash-appcore/pkg/app"
)

var (
	n  = 0
	mu sync.RWMutex
)

func main() {
	app := appcore.New()

	// runs every second
	app.AddCronJob("* * * * * *", "counter", count)

	app.Run()
}

func count() {
	mu.Lock()
	defer mu.Unlock()

	n++

	fmt.Printf("counter: %d\n", n)
}
