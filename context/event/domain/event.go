package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	MaxAllowedRoleID    = 10
	MaxAllowedChannelID = 10
)

// イベントのInterfaceです
type Event interface {
	ID() UUID
	Kind() Kind
	AllowedRoleID() []DiscordID
	AllowedChannelID() []DiscordID
}

// Eventの共通の構造体です
type EventCore struct {
	id               UUID
	kind             Kind
	allowedRoleID    []DiscordID
	allowedChannelID []DiscordID
}

// イベントを生成します
func NewEvent(
	id UUID,
	kind Kind,
	allowedRoleID []DiscordID,
	allowedChannelID []DiscordID,
) (EventCore, error) {
	e := EventCore{
		id:               id,
		kind:             kind,
		allowedRoleID:    allowedRoleID,
		allowedChannelID: allowedChannelID,
	}

	if err := e.validate(); err != nil {
		return e, errors.NewError("検証に失敗しました", err)
	}

	return e, nil
}

// イベントのIDを返します
func (e EventCore) ID() UUID {
	return e.id
}

// イベントの種類を返します
func (e EventCore) Kind() Kind {
	return e.kind
}

// イベントの許可されたロールIDを返します
func (e EventCore) AllowedRoleID() []DiscordID {
	return e.allowedRoleID
}

// イベントの許可されたチャンネルIDを返します
func (e EventCore) AllowedChannelID() []DiscordID {
	return e.allowedChannelID
}

// イベントの検証を行います
func (e EventCore) validate() error {
	if len(e.allowedRoleID) > MaxAllowedRoleID {
		return errors.NewError("許可されたロールIDの最大数を超えています")
	}

	if len(e.allowedChannelID) > MaxAllowedChannelID {
		return errors.NewError("許可されたチャンネルIDの最大数を超えています")
	}

	return nil
}

// 構造体からJSONに変換します
func (e EventCore) MarshalJSON() ([]byte, error) {
	data := struct {
		ID               UUID        `json:"id"`
		Kind             Kind        `json:"kind"`
		AllowedRoleID    []DiscordID `json:"allowed_role_id"`
		AllowedChannelID []DiscordID `json:"allowed_channel_id"`
	}{
		ID:               e.id,
		Kind:             e.kind,
		AllowedRoleID:    e.allowedRoleID,
		AllowedChannelID: e.allowedChannelID,
	}

	return json.Marshal(data)
}

// JSONから構造体に変換します
func (e *EventCore) UnmarshalJSON(b []byte) error {
	data := struct {
		ID               UUID        `json:"id"`
		Kind             Kind        `json:"kind"`
		AllowedRoleID    []DiscordID `json:"allowed_role_id"`
		AllowedChannelID []DiscordID `json:"allowed_channel_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	e.id = data.ID
	e.kind = data.Kind
	e.allowedRoleID = data.AllowedRoleID
	e.allowedChannelID = data.AllowedChannelID

	return nil
}
