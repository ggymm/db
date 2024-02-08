package index

import (
	"bytes"
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
	size := 1 + 2*8
	data := make([]byte, size)

	for n := 1; n <= 8; n++ {
		data[n] = byte(math.Ceil(float64(n) / 2))
	}

	i := 2 // 索引 + 1
	begin := head + i*2
	end := size - 1
	for n := end; n >= begin; n-- {
		data[n] = data[n-2]
	}
	t.Logf("%v", data)
}

func Test_Shift2(t *testing.T) {
	head := 1
	size := 1 + 2*8

	shift := func(data []byte, i int) {
		start := head + (i+1)*2
		length := size - head
		for n := length; n >= start; n-- {
			data[n] = data[n-2]
		}
	}

	data := make([]byte, size)
	for n := 1; n <= 8; n++ {
		if n%2 != 0 {
			data[n] = byte(n)
			data[n+1] = byte(n)
		}
	}

	var (
		key   = 2
		child = 2
	)

	var i int
	for i < 4 {
		if key <= int(data[head+2*i]) {
			break
		}
		i++
	}

	// 我写的处理方法
	handle1 := func(data []byte) []byte {
		newData := make([]byte, len(data))
		copy(newData, data)

		shift(newData, i)
		newData[head+2*i] = byte(key)
		newData[head+2*(i+1)+1] = byte(child)
		return newData
	}

	// 原版本处理方法
	handle2 := func(data []byte) []byte {
		newData := make([]byte, len(data))
		copy(newData, data)

		nextKey := newData[head+2*i]
		newData[head+2*i] = byte(key)
		shift(newData, i+1)
		newData[head+2*(i+1)] = nextKey
		newData[head+2*(i+1)+1] = byte(child)
		return newData
	}

	t.Logf("%v", handle1(data))
	t.Logf("%v", handle2(data))
	t.Logf("compare: %v", bytes.Equal(handle1(data), handle2(data)))
}

func Test_Inf(t *testing.T) {
	t.Log(uint64(1<<63) - 1 + (1 << 63))
	t.Log(uint64(math.MaxUint64))
	t.Log(float64(math.MaxUint64))
	t.Log(uint64(math.Inf(1)))
	t.Log(math.Inf(1) > math.MaxUint64)
	t.Log(uint64(math.Inf(1)) > uint64(math.MaxUint64))
	t.Log(math.Inf(1) > float64(math.MaxUint64))

	t.Log(uint64(14263705658874422041) > uint64(math.MaxUint64))
}

func Test_Data(t *testing.T) {
	buf := make([]byte, nodeSize)

	setLeaf(buf, true)
	setKeysNum(buf, 1)
	setSibling(buf, 1)
	t.Logf("%v", buf)

	setKey(buf, 0, 99)
	setChild(buf, 0, 99)
	t.Logf("%v", buf)
}
