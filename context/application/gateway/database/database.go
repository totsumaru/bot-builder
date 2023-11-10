package database

// アプリケーションのDBスキーマです
type Application struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	ServerID string `gorm:"type:uuid;not null;index:idx_server_id"`
	Data     []byte `gorm:"type:jsonb"`
}
