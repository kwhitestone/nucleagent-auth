// Package apikey 提供 API Key 管理 addon。
//
// API Key 用于 S2S（如外部系统调 core）或程序化访问，哈希存储（明文只在创建时返回一次）。
// 用户通过 /api/v1/addons/auth/api-keys 管理自己的 key。
package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// APIKey API Key 表。
//
// KeyHash 存 SHA256(明文 key)，明文只在创建时返回一次，丢失需重新生成。
// Prefix 是明文 key 的前 8 字符，用于展示识别（如 nak_xxxxxxxx...）。
type APIKey struct {
	ID        uint           `json:"id" gorm:"primarykey;comment:主键ID"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	UserID   uint   `json:"userId" gorm:"column:user_id;index;comment:所属用户ID"`
	Name     string `json:"name" gorm:"column:name;size:128;comment:Key用途说明"`
	Prefix   string `json:"prefix" gorm:"column:prefix;size:16;index;comment:明文key前缀(识别用)"`
	KeyHash  string `json:"-" gorm:"column:key_hash;size:64;uniqueIndex;comment:SHA256(明文key)"`
	LastUsed *time.Time `json:"lastUsed" gorm:"column:last_used;comment:最后使用时间"`
	Enable   bool   `json:"enable" gorm:"column:enable;default:true;comment:是否启用"`
}

// TableName 固定表名。
func (APIKey) TableName() string { return "api_keys" }

// KeyPrefix 明文 key 的前缀，便于识别。
const KeyPrefix = "nak_"

// GenerateKey 生成明文 key（nak_ + 32 字节随机 hex）。返回 (明文, 前缀, hash, error)。
func GenerateKey() (plain, prefix, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", "", err
	}
	plain = KeyPrefix + hex.EncodeToString(buf)
	prefix = plain[:len(KeyPrefix)+8] // nak_ + 8 hex
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])
	return plain, prefix, hash, nil
}

// HashKey 计算明文 key 的 SHA256 hex。
func HashKey(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
