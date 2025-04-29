package machine

import (
	"context"
	"fmt"
	"sync"
)

type Allocator interface {
	// New an allocator
	New(addr string, instanceID string)
	// Name returns the name of allocator
	Name() string
	// Acquire a machine_id
	Acquire(ctx context.Context, tenantID string) (uint16, error)
	// Release a machine_id
	Release(ctx context.Context, tenantID string, machineID uint16) error
	// Close the allocator
	Close()
}

var driversMu sync.RWMutex
var drivers = make(map[string]Allocator)

var closeFuncs []func()

func Register(name string, allocator Allocator) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if allocator == nil {
		panic("machine: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("machine: Register called twice for driver " + name)
	}
	drivers[name] = allocator
	closeFuncs = append(closeFuncs, func() {
		allocator.Close()
	})
}

func Get(name string) (allocator Allocator, err error) {
	driversMu.RLock()
	defer driversMu.RUnlock()
	if item, ok := drivers[name]; ok {
		allocator = item
		return
	} else {
		panic(fmt.Errorf("machine: driver %s not found", name))
	}

	return
}

func Close() {
	for _, closeFunc := range closeFuncs {
		closeFunc()
	}
}
