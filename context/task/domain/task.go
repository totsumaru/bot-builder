package domain

import (
	"encoding/json"
)

// タスクです
type Task struct {
	id            UUID
	serverID      DiscordID
	applicationID UUID
	ifBlock       IfBlock
}

// タスクを生成します
func NewTask(
	id UUID,
	serverID DiscordID,
	applicationID UUID,
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
func (t Task) ID() UUID {
	return t.id
}

// サーバーIDを返します
func (t Task) ServerID() DiscordID {
	return t.serverID
}

// アプリケーションIDを返します
func (t Task) ApplicationID() UUID {
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
		ID            UUID      `json:"id"`
		ServerID      DiscordID `json:"server_id"`
		ApplicationID UUID      `json:"application_id"`
		IfBlock       IfBlock   `json:"if_block"`
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
		ID            UUID      `json:"id"`
		ServerID      DiscordID `json:"server_id"`
		ApplicationID UUID      `json:"application_id"`
		IfBlock       IfBlock   `json:"if_block"`
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
