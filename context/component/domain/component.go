package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
)

// コンポーネントのinterfaceです
type Component interface {
	ID() context.UUID
	ServerID() context.DiscordID
	ApplicationID() context.UUID
}

// コンポーネントの共通の構造体です
type ComponentCore struct {
	id            context.UUID
	serverID      context.DiscordID
	applicationID context.UUID
}

// コンポーネントの共通の構造体を作成します
func NewComponentCore(
	id context.UUID,
	serverID context.DiscordID,
	applicationID context.UUID,
) (ComponentCore, error) {
	res := ComponentCore{
		id:            id,
		serverID:      serverID,
		applicationID: applicationID,
	}

	if err := res.validate(); err != nil {
		return res, err
	}

	return res, nil
}

// IDを取得します
func (c ComponentCore) ID() context.UUID {
	return c.id
}

// サーバーIDを取得します
func (c ComponentCore) ServerID() context.DiscordID {
	return c.serverID
}

// アプリケーションIDを取得します
func (c ComponentCore) ApplicationID() context.UUID {
	return c.applicationID
}

// コンポーネントの共通の構造体を検証します
func (c ComponentCore) validate() error {
	return nil
}

// コンポーネントの共通の構造体をJSONに変換します
func (c ComponentCore) MarshalJSON() ([]byte, error) {
	data := struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
	}{
		ID:            c.id,
		ServerID:      c.serverID,
		ApplicationID: c.applicationID,
	}

	return json.Marshal(data)
}

// JSONからコンポーネントの共通の構造体を復元します
func (c *ComponentCore) UnmarshalJSON(b []byte) error {
	var data struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	c.id = data.ID
	c.serverID = data.ServerID
	c.applicationID = data.ApplicationID

	return nil
}
