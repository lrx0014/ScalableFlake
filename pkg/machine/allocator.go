package machine

import "context"

type Allocator interface {
	Acquire(ctx context.Context, tenantID string) (uint64, error)
	Release(ctx context.Context, tenantID string, machineID uint64) error
}
