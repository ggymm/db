package ver

import "db/internal/tx"

// 事务缓存
//
// 用于处理事务隔离级别
// level0：readCommitted（读提交）
// level1：repeatableRead（可重复度）

type txCache struct {
	Err         error
	Tid         uint64
	Level       int
	AutoAborted bool

	snapshot map[uint64]bool // 当前事务的快照数据
}

// newTxCache
// tid：事务ID
// level：隔离级别
// active：正在执行的事务
func newTxCache(tid uint64, level int, active map[uint64]*txCache) *txCache {
	tc := new(txCache)
	tc.Tid = tid
	tc.Level = level
	if level > 0 { // 隔离级别
		tc.snapshot = make(map[uint64]bool)
		for id := range active {
			tc.snapshot[id] = true
		}
	}
	return tc
}

func (t *txCache) InSnapshot(tid uint64) bool {
	if tid == 0 {
		return false
	}
	_, ok := t.snapshot[tid]
	return ok
}

func (t *txCache) IsEntryVisible(tm tx.Manage, ent *entry) bool {
	if t.Level == 0 {
		return t.ReadCommitted(tm, ent)
	} else {
		return t.RepeatableRead(tm, ent)
	}
}

// IsEntryVersionSkip 判断是否发生了版本跳跃
func (t *txCache) IsEntryVersionSkip(tm tx.Manage, ent *entry) bool {
	if t.Level == 0 {
		return false
	} else {
		tid := t.Tid
		maxTid := ent.Max()
		if maxTid != tid && (maxTid > tid || t.InSnapshot(maxTid)) {
			return true
		}
	}
	return false
}

// ReadCommitted 读提交事务隔离级别可见性判断
//
// 读提交是指一个事务只能读取已经提交的事务产生的数据
// 判断逻辑：
// 1. 事务自身创建的数据，且未被删除
// 2. 事务 min 已经提交
// 2.1 该数据未被删除 max == 0
// 2.2 未被其他事务删除并提交 max != tid && max no committed
func (t *txCache) ReadCommitted(tm tx.Manage, ent *entry) bool {
	tid := t.Tid
	minTid := ent.Min()
	maxTid := ent.Max()

	if minTid == tid && maxTid == 0 {
		return true
	}
	if tm.IsCommitted(minTid) {
		if maxTid == 0 {
			return true
		}
		if maxTid != tid && !tm.IsCommitted(maxTid) {
			return true
		}
	}
	return false
}

// RepeatableRead 可重复读事务隔离级别可见性判断
//
// 可重复读是指一个事务只能读取已经提交的事务产生的数据，且事务开始后，其他事务产生的数据对该事务不可见
//
// 使用 snapshot 保存当前事务开始时，处于活跃状态的事务ID
// 判断逻辑：
// 1. 事务自身创建的数据，且未被删除
// 2. 事务 min 已经提交，并且在当前事务之前 min < tid && min no in snapshot
// 2.1 该数据未被删除 max == 0
// 2.2 被其他事务删除 max != tid
// 2.2.1 其他事务未提交 max no committed
// 2.2.2 其他事务在当前事务之后才开始 max > tid
// 2.2.3 其他事务在当前事务开始时还没有提交 max in snapshot
//
// 首先保证无法读取，已经被修改（更新或者删除）的数据（2的描述）
// 其次保证可以读取，已经被删除的数据（2.2的描述）
func (t *txCache) RepeatableRead(tm tx.Manage, ent *entry) bool {
	tid := t.Tid
	minTid := ent.Min()
	maxTid := ent.Max()

	if minTid == tid && maxTid == 0 {
		return true
	}
	if tm.IsCommitted(minTid) && minTid < tid && !t.InSnapshot(minTid) {
		if maxTid == 0 {
			return true
		}
		if maxTid != tid {
			if !tm.IsCommitted(maxTid) || maxTid > tid || t.InSnapshot(maxTid) {
				return true
			}
		}
	}
	return false
}
