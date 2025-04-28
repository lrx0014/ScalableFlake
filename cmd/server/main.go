package main

import (
	"os"

	"github.com/lrx0014/ScalableFlake/internal/server"
	"github.com/lrx0014/ScalableFlake/pkg/driver"
)

func main() {
	backend := os.Getenv("UID_BACKEND") // "redis" or "etcd"
	server.RunServer(driver.GetDriver(backend))
}
