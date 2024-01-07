package cache

import "testing"

func TestNewCache(t *testing.T) {
	obtain := func(key uint64) (any, error) {
		return key, nil
	}
	release := func(key any) {

	}

	ops := new(Option)
	ops.Obtain = obtain
	ops.Release = release
	ops.MaxCount = 10
	c := NewCache(ops)

	t.Logf("%+v", c)
}

func TestCache_Batch(t *testing.T) {
}
