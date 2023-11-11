package database

// タスクのDBスキーマです
type Task struct {
	ID            string `gorm:"type:uuid;primary_key;"`
	ServerID      string `gorm:"type:varchar(255);not null;index:idx_server_id"`
	ApplicationID string `gorm:"type:uuid;not null;index:idx_application_id"`
	Data          []byte `gorm:"type:jsonb"`
}
