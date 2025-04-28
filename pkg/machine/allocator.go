package machine

import "context"

type Allocator interface {
	// Acquire a machine_id
	Acquire(ctx context.Context, tenantID string) (uint16, error)
	// Release a machine_id
	Release(ctx context.Context, tenantID string, machineID uint64) error
}
