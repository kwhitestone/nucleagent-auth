package apikey

import (
	"errors"
	"time"

	"whitestone.top/prism-fusion/global"

	"gorm.io/gorm"
)

// Service API Key 业务逻辑。
type Service struct{}

// Create 为 userID 创建一个新 key，返回含明文的视图（明文只此一次返回）。
func (s *Service) Create(userID uint, name string) (*APIKeyView, error) {
	plain, prefix, hash, err := GenerateKey()
	if err != nil {
		return nil, err
	}
	key := &APIKey{
		UserID:  userID,
		Name:    name,
		Prefix:  prefix,
		KeyHash: hash,
		Enable:  true,
	}
	if err := global.PRISM_DB.Create(key).Error; err != nil {
		return nil, err
	}
	return &APIKeyView{
		ID:        key.ID,
		Name:      key.Name,
		Prefix:    key.Prefix,
		Plaintext: plain,
		Enable:    key.Enable,
		CreatedAt: key.CreatedAt,
	}, nil
}

// List 列出用户的所有 key（不含明文/ hash）。
func (s *Service) List(userID uint) ([]APIKeyView, error) {
	var keys []APIKey
	if err := global.PRISM_DB.Where("user_id = ?", userID).
		Order("id DESC").Find(&keys).Error; err != nil {
		return nil, err
	}
	views := make([]APIKeyView, 0, len(keys))
	for i := range keys {
		views = append(views, APIKeyView{
			ID:        keys[i].ID,
			Name:      keys[i].Name,
			Prefix:    keys[i].Prefix,
			Enable:    keys[i].Enable,
			LastUsed:  keys[i].LastUsed,
			CreatedAt: keys[i].CreatedAt,
		})
	}
	return views, nil
}

// Delete 软删除用户的一个 key。
func (s *Service) Delete(userID, id uint) error {
	res := global.PRISM_DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&APIKey{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Verify 校验明文 key，返回所属 userID。命中时更新 LastUsed。
// 未启用或不存在返回 ErrInvalidKey。
func (s *Service) Verify(plain string) (uint, error) {
	if len(plain) < len(KeyPrefix) || plain[:len(KeyPrefix)] != KeyPrefix {
		return 0, ErrInvalidKey
	}
	hash := HashKey(plain)
	var key APIKey
	if err := global.PRISM_DB.Where("key_hash = ? AND enable = ?", hash, true).
		First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrInvalidKey
		}
		return 0, err
	}
	now := time.Now()
	global.PRISM_DB.Model(&APIKey{}).Where("id = ?", key.ID).
		Update("last_used", now)
	return key.UserID, nil
}

// ErrInvalidKey key 无效或已停用。
var ErrInvalidKey = errors.New("invalid or disabled api key")

// APIKeyView API Key 的对外视图（创建时带 Plaintext，其余场景不带）。
type APIKeyView struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	Plaintext string     `json:"plaintext,omitempty"` // 仅创建时返回一次
	Enable    bool       `json:"enable"`
	LastUsed  *time.Time `json:"lastUsed,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}
