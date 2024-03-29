package ver

// 实际场景描述
//
// t1 begin
// t2 begin
// t1 read x
// t2 read x
// t1 update x = 1
// t2 update x = 2
// t1 commit
// t2 commit
//
// 两段锁协议 会在 t1 update x+1 时加锁，阻塞 t2 的任何操作

// 不可重复读场景描述
//
// 场景1
// t1 begin
// t1 read x
// t2 begin
// t2 update x = 1
// t2 commit
// t1 read x
//
// 此时 entry 内容为
// t1 begin --- t1 read --- t2 begin
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t0        |       0        |       0        |
// +----------------+----------------+----------------+
// t2 update --- t2 commit --- t1 read
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t2        |       0        |       1        |
// +----------------+----------------+----------------+
// 结论，t1 无法读取到当前数据
//
// 场景2
// t1 begin
// t2 begin
// t2 read x
// t1 delete
// t1 commit
// t2 read x
//
// 此时 entry 内容为
// t1 begin --- t2 begin --- t2 read
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t0        |       0        |       0        |
// +----------------+----------------+----------------+
// t1 delete --- t1 commit --- t2 read
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t0        |      t1        |       0        |
// +----------------+----------------+----------------+
// 结论，t2 可以读取到当前数据（即虽然删除了，但是当前数据对于 t2 还是可见的）
//
// 解决此问题
// 1. t1 不能读取到 t1 开始后的数据
// 1. t1 不能读取到 t1 开始时还处于 active 的数据

// 版本跳跃场景描述
//
// 场景1
// t1 begin
// t2 begin
// t1 read x
// t2 read x
// t1 update x = 1
// t1 commit
// t2 update x = 2
// t2 commit
//
// 此时 entry 内容为
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t1        |       0        |       1        |
// +----------------+----------------+----------------+
//
// 场景2
// t1 begin
// t2 begin
// t2 read x
// t1 delete x
// t1 commit
// t2 update x = 2
// t2 commit
//
// 此时 entry 内容为
// +----------------+----------------+----------------+
// |      min       |      max       |      data      |
// +----------------+----------------+----------------+
// |      t1        |      t1        |       1        |
// +----------------+----------------+----------------+
//
// 解决此问题
// 1. max（t1） 必须已经提交
// 2. max（t1） > t2 （非此例情况）或者 max（t1） 在 t2 开始前未提交
