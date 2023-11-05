package database

// Actionのスキーマです
//
// 構造体名がDBのテーブル名になります。
type Action struct {
	ID      string `gorm:"type:uuid;primary_key;"`
	EventID string `gorm:"type:uuid;not null;"`
	Data    []byte `gorm:"type:jsonb"`
}
