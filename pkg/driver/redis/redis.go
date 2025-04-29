package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/lrx0014/ScalableFlake/internal/shutdown"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

var driverName = "redis"

func init() {
	allocator.Register(driverName, &AllocatorRedis{})
}

const (
	KeyPrefix = "scalableflake:machine_id:%s:%d"
)

type AllocatorRedis struct {
	client     *redis.Client
	maxID      int
	ttl        time.Duration
	instanceID string
	ownedKeys  []string
}

func (r *AllocatorRedis) New(addr, instanceID string) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	if pong != "PONG" {
		panic("redis server is not connected")
	}

	r.client = rdb
	r.instanceID = instanceID
	r.maxID = 1023
	r.ttl = 5 * time.Minute
	return
}

func (r *AllocatorRedis) Name() string {
	return driverName
}

func (r *AllocatorRedis) Acquire(ctx context.Context, tenantID string) (machineId uint16, err error) {
	for i := 0; i <= r.maxID; i++ {
		key := fmt.Sprintf(KeyPrefix, tenantID, i)
		var success bool
		success, err = r.client.SetNX(ctx, key, r.instanceID, r.ttl).Result()
		if err != nil {
			log.Errorf("redis setnx err: %v", err)
			return
		}
		if success {
			r.ownedKeys = append(r.ownedKeys, key)
			machineId = uint16(i)
			log.Infof("instance %s: acquired machine id: %d for tenant %s", r.instanceID, machineId, tenantID)
			r.startLeaseRenewal(ctx, tenantID, machineId, 1*time.Minute)
			return
		}
	}

	panic(fmt.Errorf("no available machine ID for tenant %s", tenantID))
	return
}

func (r *AllocatorRedis) Release(ctx context.Context, tenantID string, machineID uint16) (err error) {
	key := fmt.Sprintf(KeyPrefix, tenantID, machineID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("redis get err: %v", err)
		return err
	}

	if val == r.instanceID {
		return r.client.Del(ctx, key).Err()
	}

	return
}

func (r *AllocatorRedis) Close() {
	// release all tenant keys and machine ids
	for _, key := range r.ownedKeys {
		ss := strings.Split(key, ":")
		// scalableflake:machine_id:%s:%d
		if len(ss) == 4 && ss[0] == "scalableflake" && ss[1] == "machine_id" {
			tenantID := ss[2]
			machineId, e := strconv.Atoi(ss[3])
			if e == nil {
				_ = r.Release(context.Background(), tenantID, uint16(machineId))
				log.Infof("instance %s: released machine id %d for tenant %s", r.instanceID, machineId, tenantID)
			}
		}
	}

	_ = r.client.Close()
	log.Infof("redis allocator closed")
}

func (r *AllocatorRedis) startLeaseRenewal(ctx context.Context, tenantID string, machineID uint16, interval time.Duration) {
	go func() {
		key := fmt.Sprintf(KeyPrefix, tenantID, machineID)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				val, err := r.client.Get(ctx, key).Result()
				if errors.Is(err, redis.Nil) {
					panic(fmt.Errorf("machine ID missing for tenant %s", tenantID))
					return
				}
				if err != nil {
					log.Errorf("redis get err: %v", err)
					continue
				}

				if val == r.instanceID {
					err = r.client.Expire(ctx, key, r.ttl).Err()
					if err != nil {
						log.Errorf("redis expire err: %v", err)
					} else {
						log.Infof("Renewed machine ID for tenant %s", tenantID)
					}
				} else {
					panic(fmt.Errorf("[WARN] MachineID key ownership lost: %s", key))
					return
				}

			case <-shutdown.Done():
				log.Infof("lease renewal stopped.")
				return
			}
		}
	}()
}
