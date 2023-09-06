package fresh

import (
	"fmt"

	"github.com/freshcloud-io/go-sdk/resource"
)

type ResourceOption func(*Resources)

type Resources struct {
	CPU    int64
	Memory int64
	Disk   int64
}

func WithResourcesCPU(value string) ResourceOption {
	return func(r *Resources) {
		qty := resource.MustParse(value)
		r.CPU = qty.MilliValue()
	}
}

func WithResourcesMemory(value string) ResourceOption {
	return func(r *Resources) {
		qty := resource.MustParse(value)
		r.Memory = qty.Value()
	}
}

func WithResourcesDisk(value string) ResourceOption {
	return func(r *Resources) {
		qty := resource.MustParse(value)
		r.Disk = qty.Value()
	}
}

func NewResource(opts ...ResourceOption) *Resources {
	const (
		// TODO: check if these are correct
		defaultCPU    = 250
		defaultMemory = 128
		defaultDisk   = 1000
	)

	r := &Resources{
		CPU:    defaultCPU,
		Memory: defaultMemory,
		Disk:   defaultDisk,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Resources) String() string {
	return fmt.Sprintf("cpu: %d\tmemory: %d\tdisk: %d", r.CPU, r.Memory, r.Disk)
}
