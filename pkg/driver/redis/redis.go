package redis

import (
	"context"
	"fmt"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type AllocatorRedis struct {
	client     *redis.Client
	keyPrefix  string
	lockPrefix string
	maxID      int64
	lockTTL    time.Duration
}

func NewRedisAllocator(addr string) allocator.Allocator {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return &AllocatorRedis{
		client:     client,
		keyPrefix:  "uid_gen:",
		lockPrefix: "uid_lock:",
		maxID:      1023,
		lockTTL:    5 * time.Second,
	}
}

func (r *AllocatorRedis) Acquire(ctx context.Context, tenantID string) (machineId uint16, err error) {
	if tenantID == "" {
		tenantID = "default"
	}

	log.Infof("init: tenant: %s", tenantID)
	log.Infof("init: driver: redis")

	defer func() {
		log.Infof("acquired machine id: %d", machineId)
	}()

	lockKey := r.lockPrefix + tenantID
	counterKey := r.keyPrefix + tenantID

	success, err := r.client.SetNX(ctx, lockKey, "1", r.lockTTL).Result()
	if err != nil {
		return 0, fmt.Errorf("redis setnx lock error: %w", err)
	}
	if !success {
		return 0, fmt.Errorf("cannot acquire lock for tenant %s, please retry", tenantID)
	}
	defer func() {
		r.client.Del(ctx, lockKey)
	}()

	id, err := r.client.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr error: %w", err)
	}

	if id <= r.maxID {
		machineId = uint16(id)
		return
	}

	script := `
        redis.call("SET", KEYS[1], 0)
        local val = redis.call("INCR", KEYS[1])
        return val
    `
	id, err = r.client.Eval(ctx, script, []string{counterKey}).Int64()
	if err != nil {
		return 0, fmt.Errorf("redis reset script error: %w", err)
	}

	log.Infof("reset machine id due to reaching maxinum")
	machineId = uint16(id)

	return
}

func (r *AllocatorRedis) Release(ctx context.Context, tenantID string, machineID uint64) error {
	return nil
}
