package index

import (
	"math"
	"testing"
)

func Test_Ceil(t *testing.T) {
	t.Log(float64(1 / 2))
	t.Log(math.Ceil(float64(1 / 2)))

	t.Log(float64(1) / 2)
	t.Log(math.Ceil(float64(1) / 2))
}

func Test_Shift(t *testing.T) {
	head := 1
	size := 1 + 2*10
	data := make([]byte, size)

	for n := 1; n <= 8; n++ {
		data[n] = byte(math.Ceil(float64(n) / 2))
	}

	i := 3
	begin := head + i*2
	end := size - 1
	for n := end; n >= begin; n-- {
		data[n] = data[n-2]
	}
	t.Logf("%v", data)
}
