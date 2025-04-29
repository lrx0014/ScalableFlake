package main

import (
	"github.com/lrx0014/ScalableFlake/internal/server"
	"github.com/lrx0014/ScalableFlake/internal/shutdown"
	"github.com/lrx0014/ScalableFlake/pkg/snowflake"
	log "github.com/sirupsen/logrus"

	// the underlying database
	_ "github.com/lrx0014/ScalableFlake/pkg/driver/redis"
)

func main() {
	shutdown.Init(func() {
		snowflake.Close()
	})

	go server.RunServer()

	<-shutdown.Done()
	log.Infof("ScalableFlake stopped")
}
