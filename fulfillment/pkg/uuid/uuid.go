// Package uuid 集中生成主键：同一个「唯一 ID」，按聚合主键的类型给两种形态。
//
//	NewString() → string 形式（36 位小写 UUID）
//	NewInt64()  → int64 形式（UUID 折成 64 位正整数）
//
// 各聚合不要再自己写 generateID / uuid.New().String()，统一走本包 ——
// 将来要换生成策略（雪花 ID / 数据库自增）只改这一个地方。
package uuid

import (
	"encoding/binary"

	guuid "github.com/google/uuid"
)

// NewString 生成 string 形式的主键：36 位小写 UUID（带连字符）。
func NewString() string {
	return guuid.NewString()
}

// NewInt64 生成 int64 形式的主键。
func NewInt64() int64 {
	return Int64(guuid.New())
}

// Int64 把 UUID 折成 int64：取前 8 字节按大端解释，最高位清零保证是正数。
//
// ⚠️ 128 位压成 63 位后随机性只剩 63 bit —— 按生日界，大约 2^31 个 ID 时撞一次的概率就到 50%。
// 主键量级大、或要求强唯一时，把本函数换成雪花 ID / 数据库自增即可（调用方不用改）。
func Int64(u guuid.UUID) int64 {
	return int64(binary.BigEndian.Uint64(u[:8]) & 0x7FFFFFFFFFFFFFFF)
}

// StringOf 把 UUID 渲染成 string 形式（与 NewString 同一口径）。
func StringOf(u guuid.UUID) string {
	return u.String()
}
