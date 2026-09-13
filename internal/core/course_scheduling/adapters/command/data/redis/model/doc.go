// Package model 是 course_scheduling 写侧 Redis 的序列化模型。
//
// 与 mysql/model 一一对应：每个聚合一个文件，struct 字段与 MySQL 侧一致
// （值对象同样就地打平），区别只在两处：
//
//  1. 用 json tag 而不是数据库列注释 —— Redis 存的是字节，序列化走 JSON。
//  2. 转换函数叫 XxxToRedis / XxxFromRedis（MySQL 侧是 XxxToPO / XxxToDO）。
//
// 它是存储形态，不是领域模型：领域对象（domain/aggregate/*）不带任何序列化
// tag，转换只发生在这一层，跟 MySQL 的 PO 是同一个道理。
//
// 快照里不含 MySQL 的乐观锁列 LockVersion —— 那是写库并发控制的东西，
// 和「缓存一份聚合」无关。
package model
