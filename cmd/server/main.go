package main

import (
	"github.com/lrx0014/ScalableFlake/internal/server"
	"github.com/lrx0014/ScalableFlake/pkg/snowflake"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go server.RunServer()
	waitForShutdown(func() {
		snowflake.Close()
	})
}

func waitForShutdown(releaseFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Infof("Graceful shutdown initiated...")
	releaseFunc()
}
