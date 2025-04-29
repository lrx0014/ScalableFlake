package snowflake

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/lrx0014/ScalableFlake/internal/env"
	"github.com/lrx0014/ScalableFlake/pkg/machine"
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"os"
	"time"

	_ "github.com/lrx0014/ScalableFlake/pkg/driver/redis"
)

var sfs map[string]*sonyflake.Sonyflake
var allocators []machine.Allocator

func init() {
	sfs = map[string]*sonyflake.Sonyflake{}
	backend, tenant, addr := config()
	allocator, err := machine.Get(backend)
	if err != nil {
		panic(err)
	}
	allocator.New(addr, generateInstanceID())
	log.Infof("allocator with backend [%s] initiated", allocator.Name())
	allocators = append(allocators, allocator)
	settings := sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return allocator.Acquire(context.Background(), tenant)
		},
	}
	sf := sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("sonyflake not created")
	}

	// initial a default tenant
	sfs[tenant] = sf
}

func GenerateUID(tenant string) (uid uint64, err error) {
	sf := sfs[tenant]
	if sf == nil {
		err = fmt.Errorf("tenant %s not exist", tenant)
		return
	}

	uid, err = sf.NextID()
	return
}

func Close() {
	for _, allocator := range allocators {
		allocator.Close()
	}
	log.Infof("closed sonyflake")
}

func config() (backend, tenant, addr string) {
	backend = os.Getenv(env.BACKEND) // "redis" or "etcd"
	tenant = os.Getenv(env.TENANT)
	addr = os.Getenv(env.ADDR)

	if backend == "" {
		backend = "redis"
	}
	if tenant == "" {
		tenant = "default"
	}
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	return
}

func generateInstanceID() string {
	b := make([]byte, 5) // 10 hex characters
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random instance ID")
	}
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("%s-%d", hex.EncodeToString(b), timestamp)
}
