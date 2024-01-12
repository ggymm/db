package page

import "testing"

func TestNewIndex(t *testing.T) {
	pi := NewIndex()
	t.Logf("%+v", pi)
}
