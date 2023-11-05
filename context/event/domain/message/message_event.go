package message

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/event/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// メッセージイベントの構造体です
type MessageEvent struct {
	domain.Event
	keyword   Keyword
	matchType MatchType
}

// メッセージイベントを生成します
func NewMessageEvent(
	allowedRoleID []domain.DiscordID,
	allowedChannelID []domain.DiscordID,
	keyword Keyword,
	matchType MatchType,
) (MessageEvent, error) {
	id, err := domain.NewUUID()
	if err != nil {
		return MessageEvent{}, errors.NewError("UUIDを生成できません", err)
	}

	k, err := domain.NewKind(domain.EventKindMessageCreate)
	if err != nil {
		return MessageEvent{}, errors.NewError("イベントの種類を生成できません", err)
	}

	e, err := domain.NewEvent(id, k, allowedRoleID, allowedChannelID)
	if err != nil {
		return MessageEvent{}, errors.NewError("イベントを作成できません", err)
	}

	me := MessageEvent{
		Event:     e,
		keyword:   keyword,
		matchType: matchType,
	}

	if err = me.Validate(); err != nil {
		return me, errors.NewError("検証に失敗しました", err)
	}

	return me, nil
}

// イベントのキーワードを返します
func (e MessageEvent) Keyword() Keyword {
	return e.keyword
}

// イベントの一致条件を返します
func (e MessageEvent) MatchType() MatchType {
	return e.matchType
}

// 検証します
func (e MessageEvent) Validate() error {
	return nil
}

// 構造体からJSONに変換します
func (e MessageEvent) MarshalJSON() ([]byte, error) {
	data := struct {
		ID               domain.UUID        `json:"id"`
		Kind             domain.Kind        `json:"kind"`
		AllowedRoleID    []domain.DiscordID `json:"allowed_role_id"`
		AllowedChannelID []domain.DiscordID `json:"allowed_channel_id"`
		Keyword          Keyword            `json:"keyword"`
		MatchType        MatchType          `json:"match_type"`
	}{
		ID:               e.ID(),
		Kind:             e.Kind(),
		AllowedRoleID:    e.AllowedRoleID(),
		AllowedChannelID: e.AllowedChannelID(),
		Keyword:          e.Keyword(),
		MatchType:        e.MatchType(),
	}

	return json.Marshal(data)
}

// JSONから構造体に変換します
func (e *MessageEvent) UnmarshalJSON(b []byte) error {
	data := struct {
		ID               domain.UUID        `json:"id"`
		Kind             domain.Kind        `json:"kind"`
		AllowedRoleID    []domain.DiscordID `json:"allowed_role_id"`
		AllowedChannelID []domain.DiscordID `json:"allowed_channel_id"`
		Keyword          Keyword            `json:"keyword"`
		MatchType        MatchType          `json:"match_type"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	eventData, err := domain.NewEvent(
		data.ID,
		data.Kind,
		data.AllowedRoleID,
		data.AllowedChannelID,
	)
	if err != nil {
		return err
	}

	e.Event = eventData
	e.keyword = data.Keyword
	e.matchType = data.MatchType

	return nil
}
