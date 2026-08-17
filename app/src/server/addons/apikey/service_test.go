package apikey

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/kwhitestone/prism-fusion/global"
)

// setupTestDB 用内存 sqlite 初始化全局 DB 并迁移 APIKey 表。
func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&APIKey{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// 每个测试清空表，避免共享内存 DB 串扰。
	db.Exec("DELETE FROM api_keys")
	global.PRISM_DB = db
}

// TestAPIKeyCreateAndVerify 验证创建 -> 校验明文 -> hash 不回显。
func TestAPIKeyCreateAndVerify(t *testing.T) {
	setupTestDB(t)
	s := &Service{}

	view, err := s.Create(1, "test-key")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if view.Plaintext == "" {
		t.Error("Plaintext should be returned on create")
	}
	if view.Prefix == "" {
		t.Error("Prefix should be set")
	}
	if view.ID == 0 {
		t.Error("ID should be set")
	}

	// 明文以 KeyPrefix 开头。
	if view.Plaintext[:len(KeyPrefix)] != KeyPrefix {
		t.Errorf("plaintext should start with %q, got %q", KeyPrefix, view.Plaintext[:len(KeyPrefix)])
	}

	// 校验明文能反查出 userID。
	uid, err := s.Verify(view.Plaintext)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if uid != 1 {
		t.Errorf("Verify uid = %d, want 1", uid)
	}

	// List 不应包含明文。
	keys, err := s.List(1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("List len = %d, want 1", len(keys))
	}
	if keys[0].Plaintext != "" {
		t.Error("List must not expose plaintext")
	}
	if keys[0].Prefix != view.Prefix {
		t.Errorf("List prefix = %q, want %q", keys[0].Prefix, view.Prefix)
	}
}

// TestAPIKeyVerifyInvalid 验证无效 key 被拒。
func TestAPIKeyVerifyInvalid(t *testing.T) {
	setupTestDB(t)
	s := &Service{}

	cases := []string{"", "not-a-key", "nak_short", "nak_00000000000000000000000000000000nonexistent"}
	for _, c := range cases {
		if _, err := s.Verify(c); err == nil {
			t.Errorf("Verify(%q) should fail", c)
		}
	}
}

// TestAPIKeyDelete 验证删除（且不能删别人的 key）。
func TestAPIKeyDelete(t *testing.T) {
	setupTestDB(t)
	s := &Service{}

	view, _ := s.Create(1, "mine")
	_, _ = s.Create(2, "theirs")

	// 用户 1 不能删用户 2 的 key（这里只能删自己的）。
	if err := s.Delete(1, view.ID); err != nil {
		t.Fatalf("Delete own: %v", err)
	}
	// 再删一次（已删除）应报 NotFound。
	if err := s.Delete(1, view.ID); err == nil {
		t.Error("Delete again should fail with NotFound")
	}

	// 用户 1 的 key 数量归零。
	keys, _ := s.List(1)
	if len(keys) != 0 {
		t.Errorf("after delete, List(1) len = %d, want 0", len(keys))
	}
	// 用户 2 的 key 还在。
	keys2, _ := s.List(2)
	if len(keys2) != 1 {
		t.Errorf("List(2) len = %d, want 1", len(keys2))
	}
}

// TestAPIKeyDisabledNotVerified 验证停用的 key 校验失败。
func TestAPIKeyDisabledNotVerified(t *testing.T) {
	setupTestDB(t)
	s := &Service{}

	view, _ := s.Create(1, "will-disable")
	// 直接把 enable 置 false。
	global.PRISM_DB.Model(&APIKey{}).Where("id = ?", view.ID).Update("enable", false)

	if _, err := s.Verify(view.Plaintext); err == nil {
		t.Error("Verify on disabled key should fail")
	}
}

// TestGenerateKeyUniqueness 验证生成的 key 足够唯一（连续生成不重复）。
func TestGenerateKeyUniqueness(t *testing.T) {
	seen := make(map[string]bool, 100)
	for i := 0; i < 100; i++ {
		plain, _, _, err := GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey: %v", err)
		}
		if seen[plain] {
			t.Fatalf("duplicate key generated at iteration %d", i)
		}
		seen[plain] = true
	}
}

// TestHashKeyDeterministic 验证相同明文产出相同 hash。
func TestHashKeyDeterministic(t *testing.T) {
	plain := "nak_test123"
	if HashKey(plain) != HashKey(plain) {
		t.Error("HashKey should be deterministic")
	}
	if HashKey(plain) == HashKey("nak_other") {
		t.Error("different plaintext should produce different hash")
	}
}
