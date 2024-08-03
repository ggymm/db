参考文章：
https://juejin.cn/post/6844903808376504327
https://www.cnblogs.com/jelly12345/p/14889331.html
https://www.cnblogs.com/bytesfly/p/mysql-transaction-innodb-mvcc.html

## 事务

### 事务属性

原子性（Atomicity）

一致性（Consistency）

隔离性（Isolation）

持久性（Durability）

### 自动提交

事务被设计成自动提交的

这意味着如果你不显式地开始一个事务，每个查询都将被当做一个事务来自动创建并提交。

### 并发问题

#### 脏读（Dirty Read）

描述： 一个事务读取到另一个事务未提交的数据

场景：

1. T1 begin
2. T1 read x
3. T2 begin
4. T2 update x = 100
5. T1 read x （**dirty read**）
6. T2 commit or rollback
7. T1 read x
8. T1 commit

在步骤 5 中，事务 A 读取到了事务 B 未提交的中间状态数据，产生了脏读

#### 不可重复读（Non-Repeatable Read）

描述：同一个事务内，多次读取同一数据，读取到的数据不一致。通常是因为数据被其他事务修改

场景：

1. T1 begin
2. T1 read x
3. T2 begin
4. T2 update x = 100
5. T1 read x （**non-repeatable read**）
6. T2 commit
7. T1 read x （**non-repeatable read**）
8. T1 commit

在步骤 5 中，事务 A 没有读取到事务 B 未提交的中间状态，没有产生脏读

但是在步骤 6 中，事务 B 提交事务后，事务 A 读取到了事务 B 修改的数据

对于事务 A 两次读取的数据不一致，产生了不可重复读

#### 幻读（Phantom Read）

描述：同一个事务内，多次查询结果不一致。通常是因为数据被其他事务插入或删除

场景：

1. T1 begin
2. T1 read x > 100
3. T2 begin
4. T2 insert x = 200
5. T2 commit
6. T1 read x > 100 （**phantom read**）
7. T1 commit

事务 A 再次查询时，查询到了事务 B 插入或修改的数据，产生了幻读

幻读问题通常需要在读取时加锁（锁表），锁定后续事务的任何操作才可以解决，会影响并发性能

### 事务隔离级别

#### 读未提交（Read Uncommitted）

特性：

* 事务 A 可以读取到事务 B 未提交的数据
* 会出现**脏读**问题，应用在不严格要求数据一致性的场景

#### 读已提交（Read Committed）

特性：

* 事务 A 只能读取到事务 B 已提交的数据
* 不会出现**脏读**问题，会出现**不可重复读**问题和**幻读**问题

#### 可重复读（Repeatable Read）

特性：

* 事务 A 在同一个事务内，多次读取同一数据，读取到的数据一致
* 不会出现**脏读**问题和**不可重复读**问题，会出现**幻读**问题

#### 串行化（Serializable）

特性：

* 所有事务按照顺序独立执行
* 不会出现**脏读**问题、**不可重复读**问题和**幻读**问题
* 适用于对数据一致性要求非常高的场景，可以容忍较低并发的应用。例如：金融系统

#### 隔离级别对比

| 隔离级别 | 脏读  | 不可重复读 | 幻读  |
|------|-----|-------|-----|
| 未提交读 | 可能  | 可能    | 可能  |
| 提交读  | 不可能 | 可能    | 可能  |
| 可重复读 | 不可能 | 不可能   | 可能  |
| 序列化  | 不可能 | 不可能   | 不可能 |

当前数据库设计支持的是，**读已提交** 和 **可重复读**

### MVCC

