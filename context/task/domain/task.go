package domain

import (
	"encoding/json"
)

// タスクです
type Task struct {
	id       UUID
	serverID DiscordID
	appID    UUID
	ifBlock  IfBlock
}

// タスクを生成します
func NewTask(
	id UUID,
	serverID DiscordID,
	appID UUID,
	ifBlock IfBlock,
) (Task, error) {
	task := Task{
		id:       id,
		serverID: serverID,
		appID:    appID,
		ifBlock:  ifBlock,
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
func (t Task) AppID() UUID {
	return t.appID
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
		ID       UUID      `json:"id"`
		ServerID DiscordID `json:"server_id"`
		AppID    UUID      `json:"app_id"`
		IfBlock  IfBlock   `json:"if_block"`
	}{
		ID:       t.id,
		ServerID: t.serverID,
		AppID:    t.appID,
		IfBlock:  t.ifBlock,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (t *Task) UnmarshalJSON(data []byte) error {
	var v struct {
		ID       UUID      `json:"id"`
		ServerID DiscordID `json:"server_id"`
		AppID    UUID      `json:"app_id"`
		IfBlock  IfBlock   `json:"if_block"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.id = v.ID
	t.serverID = v.ServerID
	t.appID = v.AppID
	t.ifBlock = v.IfBlock

	return nil
}
