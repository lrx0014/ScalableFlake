package snowflake

import (
	"context"
	"errors"
	"github.com/lrx0014/ScalableFlake/pkg/driver"
	"github.com/sony/sonyflake"
	"os"
)

var sf *sonyflake.Sonyflake

func init() {
	backend := os.Getenv("UID_BACKEND") // "redis" or "etcd"
	tenant := os.Getenv("UID_TENANT")
	allocator := driver.GetDriver(backend)
	if allocator == nil {
		panic("backend error")
	}
	settings := sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return allocator.Acquire(context.Background(), tenant)
		},
	}
	sf = sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GenerateUID() (uid uint64, err error) {
	if sf == nil {
		err = errors.New("sonyflake not created")
		return
	}

	uid, err = sf.NextID()
	return
}
