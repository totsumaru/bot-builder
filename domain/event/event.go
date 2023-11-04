package event

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	MaxAllowedRoleID    = 10
	MaxAllowedChannelID = 10
)

// Eventの共通の構造体です
type Event struct {
	id               domain.UUID
	kind             Kind
	allowedRoleID    []domain.DiscordID
	allowedChannelID []domain.DiscordID
}

// イベントを生成します
func NewEvent(
	id domain.UUID,
	kind Kind,
	allowedRoleID []domain.DiscordID,
	allowedChannelID []domain.DiscordID,
) (Event, error) {
	e := Event{
		id:               id,
		kind:             kind,
		allowedRoleID:    allowedRoleID,
		allowedChannelID: allowedChannelID,
	}

	if err := e.Validate(); err != nil {
		return e, errors.NewError("検証に失敗しました", err)
	}

	return e, nil
}

// イベントのIDを返します
func (e Event) ID() domain.UUID {
	return e.id
}

// イベントの種類を返します
func (e Event) Kind() Kind {
	return e.kind
}

// イベントの許可されたロールIDを返します
func (e Event) AllowedRoleID() []domain.DiscordID {
	return e.allowedRoleID
}

// イベントの許可されたチャンネルIDを返します
func (e Event) AllowedChannelID() []domain.DiscordID {
	return e.allowedChannelID
}

// イベントの検証を行います
func (e Event) Validate() error {
	if len(e.allowedRoleID) > MaxAllowedRoleID {
		return errors.NewError("許可されたロールIDの最大数を超えています")
	}

	if len(e.allowedChannelID) > MaxAllowedChannelID {
		return errors.NewError("許可されたチャンネルIDの最大数を超えています")
	}

	return nil
}

// 構造体からJSONに変換します
func (e Event) MarshalJSON() ([]byte, error) {
	data := struct {
		ID               domain.UUID        `json:"id"`
		Kind             Kind               `json:"kind"`
		AllowedRoleID    []domain.DiscordID `json:"allowed_role_id"`
		AllowedChannelID []domain.DiscordID `json:"allowed_channel_id"`
	}{
		ID:               e.id,
		Kind:             e.kind,
		AllowedRoleID:    e.allowedRoleID,
		AllowedChannelID: e.allowedChannelID,
	}

	return json.Marshal(data)
}

// JSONから構造体に変換します
func (e *Event) UnmarshalJSON(b []byte) error {
	data := struct {
		ID               domain.UUID        `json:"id"`
		Kind             Kind               `json:"kind"`
		AllowedRoleID    []domain.DiscordID `json:"allowed_role_id"`
		AllowedChannelID []domain.DiscordID `json:"allowed_channel_id"`
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
