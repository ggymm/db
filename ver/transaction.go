package ver

import (
	"github.com/ggymm/db/tx"
)

// 事务缓存
//
// 用于处理事务隔离级别
// level0：readCommitted（读提交）
// level1：repeatableRead（可重复度）

type transaction struct {
	Err error

	Id           uint64
	Level        int
	AutoRollback bool

	snapshot map[uint64]bool // 当前事务的快照数据
}

// newTransaction
// id：事务Id
// level：隔离级别
// active：正在执行的事务
func newTransaction(id uint64, level int, active map[uint64]*transaction) *transaction {
	tc := new(transaction)
	tc.Id = id
	tc.Level = level
	if level > 0 { // 隔离级别
		tc.snapshot = make(map[uint64]bool)
		for tid := range active {
			tc.snapshot[tid] = true
		}
	}
	return tc
}

func (t *transaction) InSnapshot(tid uint64) bool {
	if tid == 0 {
		return false
	}
	_, ok := t.snapshot[tid]
	return ok
}

// IsSkip 判断是否发生了版本跳跃
// 两种情况
// 1. 被更新，此时 max 为 0
// 2. 被删除，此时 max 为
// 2.1 在当前事务之后创建的事务
// 2.2 在当前事务之前未提交的事务
func (t *transaction) IsSkip(tm tx.Manage, ent *entry) bool {
	if t.Level == 0 {
		return false
	} else {
		tid := t.Id
		maxId := ent.Max()
		if tm.IsCommitted(maxId) && (maxId > tid || t.InSnapshot(maxId)) {
			return true
		}
	}
	return false
}

func (t *transaction) IsVisible(tm tx.Manage, ent *entry) bool {
	if t.Level == 0 {
		return t.ReadCommitted(tm, ent)
	} else {
		return t.RepeatableRead(tm, ent)
	}
}

// ReadCommitted 读提交事务隔离级别可见性判断
//
// 读提交是指一个事务只能读取已经提交的事务产生的数据
// 判断逻辑：
// 1. 事务自身创建的数据，且未被删除
// 2. 事务 min 已经提交
// 2.1 该数据未被删除 max == 0
// 2.2 被其他事务删除但是没有提交 max != tid && max no committed
func (t *transaction) ReadCommitted(tm tx.Manage, ent *entry) bool {
	id := t.Id
	minId := ent.Min()
	maxId := ent.Max()

	if minId == id && maxId == 0 {
		return true
	}
	if tm.IsCommitted(minId) {
		if maxId == 0 {
			return true
		}
		if maxId != id && !tm.IsCommitted(maxId) {
			return true
		}
	}
	return false
}

// RepeatableRead 可重复读事务隔离级别可见性判断
//
// 可重复读是指一个事务只能读取已经提交的事务产生的数据，且事务开始后，其他事务产生的数据对该事务不可见
//
// 使用 snapshot 保存当前事务开始时，处于活跃状态的事务Id
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
func (t *transaction) RepeatableRead(tm tx.Manage, ent *entry) bool {
	id := t.Id
	minId := ent.Min()
	maxId := ent.Max()

	if minId == id && maxId == 0 {
		return true
	}
	if tm.IsCommitted(minId) && minId < id && !t.InSnapshot(minId) {
		if maxId == 0 {
			return true
		}
		if maxId != id {
			if !tm.IsCommitted(maxId) || maxId > id || t.InSnapshot(maxId) {
				return true
			}
		}
	}
	return false
}
