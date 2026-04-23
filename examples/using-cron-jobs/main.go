package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/mrbagir/appfr"
)

var (
	n  = 0
	mu sync.RWMutex
)

func main() {
	app := appfr.New()

	// runs every second
	app.AddCronJob("* * * * * *", "counter", count)

	app.Run()
}

func count(_ context.Context) {
	mu.Lock()
	defer mu.Unlock()

	n++

	fmt.Printf("counter: %d\n", n)
}
