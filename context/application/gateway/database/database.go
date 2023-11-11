package database

// アプリケーションのDBスキーマです
type Application struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	ServerID string `gorm:"type:varchar(255);not null;index:idx_server_id"` // Discord サーバーID用に文字列型を使用
	Data     []byte `gorm:"type:jsonb"`
}
