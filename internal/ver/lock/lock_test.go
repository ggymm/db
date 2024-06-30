package lock

import (
	"testing"
)

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
		t.Fatalf("添加 1 -> 2 成功（此时处于死锁，应该不能成功）")
		return
	}

	t.Log("success")
	print(l.String())
}

func TestLock_Add2(t *testing.T) {
	l := NewLock()

	// 添加 1……100 -> 1……100
	for i := 1; i <= 100; i++ {
		ok, ch := l.Add(uint64(i), uint64(i))
		if !ok {
			t.Fatalf("添加 %d -> %d 失败", i, i)
			return
		}
		go func() {
			<-ch
		}()
	}

	// 添加 1……99 -> 2……100
	for i := 1; i <= 99; i++ {
		ok, ch := l.Add(uint64(i), uint64(i+1))
		if !ok {
			t.Fatalf("添加 %d -> %d 失败", i, i+1)
			return
		}
		go func() {
			<-ch
		}()
	}

	// 添加 100 -> 1
	// 此时等待队列为 1 -> 2 -> 3 -> ... -> 99 -> 100
	// 添加 100 -> 1 会导致死锁
	ok, _ := l.Add(100, 1)
	if ok {
		t.Fatalf("添加 100 -> 1 成功（此时处于死锁，应该不能成功）")
		return
	}

	// 移除任意一个依赖关系，都可以解除死锁
	l.Remove(88)

	// 添加 100 -> 1
	ok, _ = l.Add(100, 1)
	if !ok {
		t.Fatalf("添加 100 -> 1 失败（此时解除死锁，应该可以成功）")
		return
	}

	t.Log("success")
	print(l.String())
}
