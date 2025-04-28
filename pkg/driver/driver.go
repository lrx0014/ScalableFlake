package driver

import (
	"github.com/lrx0014/ScalableFlake/pkg/driver/redis"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"log"
	"os"
)

func GetDriver(backend string) (al allocator.Allocator) {
	switch backend {
	case "redis":
		addr := os.Getenv("REDIS_ADDR")
		return redis.NewRedisAllocator(addr)
	default:
		log.Fatal("Unsupported backend")
	}

	return
}
