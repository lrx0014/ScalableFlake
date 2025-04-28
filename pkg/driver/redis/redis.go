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

func (r *AllocatorRedis) Acquire(ctx context.Context, tenantID string) (uint16, error) {
	if tenantID == "" {
		tenantID = "default"
	}

	log.Infof("init: tenant: %s", tenantID)
	log.Infof("init: driver: redis")

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
		return uint16(id), nil
	}

	script := `
        redis.call("SET", KEYS[1], 0)
        local val = redis.call("INCR", KEYS[1])
        return val
    `
	newID, err := r.client.Eval(ctx, script, []string{counterKey}).Int64()
	if err != nil {
		return 0, fmt.Errorf("redis reset script error: %w", err)
	}

	return uint16(newID), nil
}

func (r *AllocatorRedis) Release(ctx context.Context, tenantID string, machineID uint64) error {
	return nil
}
