package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
)

// タスクです
type Task struct {
	id            context.UUID
	serverID      context.DiscordID
	applicationID context.UUID
	ifBlock       IfBlock
}

// タスクを生成します
func NewTask(
	id context.UUID,
	serverID context.DiscordID,
	applicationID context.UUID,
	ifBlock IfBlock,
) (Task, error) {
	task := Task{
		id:            id,
		serverID:      serverID,
		applicationID: applicationID,
		ifBlock:       ifBlock,
	}

	if err := task.validate(); err != nil {
		return Task{}, err
	}

	return task, nil
}

// IDを返します
func (t Task) ID() context.UUID {
	return t.id
}

// サーバーIDを返します
func (t Task) ServerID() context.DiscordID {
	return t.serverID
}

// アプリケーションIDを返します
func (t Task) ApplicationID() context.UUID {
	return t.applicationID
}

// ifブロックを返します
func (t Task) IfBlock() IfBlock {
	return t.ifBlock
}

// 検証します
func (t Task) validate() error {
	return nil
}

// JSONに変換します
func (t Task) MarshalJSON() ([]byte, error) {
	data := struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		IfBlock       IfBlock           `json:"if_block"`
	}{
		ID:            t.id,
		ServerID:      t.serverID,
		ApplicationID: t.applicationID,
		IfBlock:       t.ifBlock,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (t *Task) UnmarshalJSON(data []byte) error {
	var v struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		IfBlock       IfBlock           `json:"if_block"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.id = v.ID
	t.serverID = v.ServerID
	t.applicationID = v.ApplicationID
	t.ifBlock = v.IfBlock

	return nil
}
