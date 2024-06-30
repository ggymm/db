package cache

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	obtain := func(key uint64) (any, error) {
		return key, nil
	}
	release := func(key any) {
	}

	opt := new(Option)
	opt.Obtain = obtain
	opt.Release = release
	opt.MaxCount = 10
	c := NewCache(opt)

	t.Logf("%+v", c)
}

func TestCache_Batch(t *testing.T) {
}
