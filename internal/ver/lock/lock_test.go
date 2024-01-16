package lock

import "testing"

func TestNewLock(t *testing.T) {
	l := NewLock()
	t.Logf("%+v", l)
}

func TestLock_Add(t *testing.T) {
	l := NewLock()

	// 添加 1 -> 1
	ok, _ := l.Add(1, 1)
	if !ok {
		t.Fatalf("添加 1 -> 1 失败")
		return
	}

	// 添加 2 -> 1
	ok, _ = l.Add(2, 1)
	if !ok {
		t.Fatalf("添加 2 -> 1 失败")
		return
	}

	t.Log("success")
	print(l.String())
}

func TestLock_Add1(t *testing.T) {
	l := NewLock()

	// 添加 1 -> 1
	ok, _ := l.Add(1, 1)
	if !ok {
		t.Fatalf("添加 1 -> 1 失败")
		return
	}

	// 添加 2 -> 2
	ok, _ = l.Add(2, 2)
	if !ok {
		t.Fatalf("添加 2 -> 2 失败")
		return
	}

	// 添加 2 -> 1
	ok, _ = l.Add(2, 1)
	if !ok {
		t.Fatalf("添加 2 -> 1 失败")
		return
	}

	// 测试死锁
	// 添加 1 -> 2
	ok, _ = l.Add(1, 2)
	if ok {
		// 此时应该死锁
		t.Fatalf("添加 1 -> 2 成功")
		return
	}

	t.Log("success")
	print(l.String())
}
