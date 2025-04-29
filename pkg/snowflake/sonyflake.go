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
	"sync"
	"time"
)

var lockSfs sync.RWMutex
var sfs map[string]*sonyflake.Sonyflake
var allocator machine.Allocator

func init() {
	sfs = map[string]*sonyflake.Sonyflake{}
	backend, addr := config()
	alloc, err := machine.Get(backend)
	if err != nil {
		panic(err)
	}
	alloc.New(addr, generateInstanceID())
	allocator = alloc
	log.Infof("allocator with backend [%s] initiated", allocator.Name())

	// init a default tenant
	tenant := "default"
	_, err = newTenant(tenant)
	if err != nil {
		panic(err)
	}
}

func GenerateUID(tenant string) (uid uint64, err error) {
	if tenant == "" {
		tenant = "default"
	}

	sf := sfs[tenant]
	// create tenant if not exist
	if sf == nil {
		log.Warnf("tenant %s not exist, creating...", tenant)
		if sf, err = newTenant(tenant); err != nil {
			log.Errorf("create tenant %s failed: %s", tenant, err.Error())
			return
		}
	}

	uid, err = sf.NextID()
	return
}

func Close() {
	machine.Close()
	log.Infof("closed sonyflake")
}

func config() (backend, addr string) {
	backend = os.Getenv(env.BACKEND) // "redis" or "etcd"
	addr = os.Getenv(env.ADDR)

	if backend == "" {
		backend = "redis"
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

func newTenant(tenant string) (sf *sonyflake.Sonyflake, err error) {
	if allocator == nil {
		panic("no allocator created")
	}
	settings := sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return allocator.Acquire(context.Background(), tenant)
		},
	}
	sf = sonyflake.NewSonyflake(settings)
	if sf == nil {
		log.Errorf("sonyflake cannot be created for tenant %s", tenant)
		return
	}

	lockSfs.Lock()
	sfs[tenant] = sf
	lockSfs.Unlock()

	log.Infof("created new tenant %s", tenant)

	return
}
