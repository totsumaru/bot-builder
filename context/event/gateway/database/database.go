package database

// Eventのスキーマです
//
// 構造体名がDBのテーブル名になります。
type Event struct {
	ID   string `gorm:"type:uuid;primary_key;"`
	Data []byte `gorm:"type:jsonb"`
}
