package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// アプリケーションです
type Application struct {
	id       context.UUID
	serverID context.DiscordID
	name     Name
}

// アプリケーションを生成します
func NewApplication(
	id context.UUID,
	serverID context.DiscordID,
	name Name,
) (Application, error) {
	a := Application{
		id:       id,
		serverID: serverID,
		name:     name,
	}

	if err := a.validate(); err != nil {
		return a, err
	}

	return a, nil
}

// アプリケーション名を変更します
func (a *Application) UpdateName(name Name) error {
	a.name = name

	if err := a.validate(); err != nil {
		return errors.NewError("アプリケーション名の変更に失敗しました", err)
	}

	return nil
}

// IDを返します
func (a Application) ID() context.UUID {
	return a.id
}

// サーバーIDを返します
func (a Application) ServerID() context.DiscordID {
	return a.serverID
}

// アプリケーションの名前を返します
func (a Application) Name() Name {
	return a.name
}

// 検証します
func (a Application) validate() error {
	return nil
}

// アプリケーションをJSONに変換します
func (a Application) MarshalJSON() ([]byte, error) {
	data := struct {
		ID       context.UUID      `json:"id"`
		ServerID context.DiscordID `json:"server_id"`
		Name     Name              `json:"name"`
	}{
		ID:       a.id,
		ServerID: a.serverID,
		Name:     a.name,
	}

	return json.Marshal(data)
}

// JSONからアプリケーションを復元します
func (a *Application) UnmarshalJSON(b []byte) error {
	data := struct {
		ID       context.UUID      `json:"id"`
		ServerID context.DiscordID `json:"server_id"`
		Name     Name              `json:"name"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	a.id = data.ID
	a.serverID = data.ServerID
	a.name = data.Name

	return nil
}
