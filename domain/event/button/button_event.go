package button

import (
	"github.com/totsumaru/bot-builder/domain"
	"github.com/totsumaru/bot-builder/domain/event"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// ボタンのイベントの構造体です
type ButtonEvent struct {
	event.Event
}

// ボタンのイベントを生成します
func NewButtonEvent(
	id domain.UUID,
	allowedRoleID []domain.DiscordID,
	allowedChannelID []domain.DiscordID,
) (ButtonEvent, error) {
	e, err := event.NewEvent(
		id,
		event.EventKindButton,
		allowedRoleID,
		allowedChannelID,
	)
	if err != nil {
		return ButtonEvent{}, errors.NewError("イベントを作成できません", err)
	}

	be := ButtonEvent{e}

	if err = be.Validate(); err != nil {
		return ButtonEvent{}, errors.NewError("検証に失敗しました", err)
	}

	return ButtonEvent{e}, nil
}

// イベントの検証を行います
func (e ButtonEvent) Validate() error {
	return nil
}

// 構造体からJSONに変換します
func (e ButtonEvent) MarshalJSON() ([]byte, error) {
	return e.Event.MarshalJSON()
}

// JSONから構造体に変換します
func (e *ButtonEvent) UnmarshalJSON(b []byte) error {
	return e.Event.UnmarshalJSON(b)
}
