package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// アクションの共通の構造体です
type Action struct {
	id      UUID
	eventID UUID
	kind    Kind
	order   Order
}

// アクションを生成します
func NewAction(eventID UUID, kind Kind, order Order) (Action, error) {
	id, err := NewUUID()
	if err != nil {
		return Action{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	a := Action{
		id:      id,
		eventID: eventID,
		kind:    kind,
		order:   order,
	}

	if err = a.Validate(); err != nil {
		return a, errors.NewError("検証に失敗しました", err)
	}

	return a, nil
}

// アクションのIDを返します
func (a Action) ID() UUID {
	return a.id
}

// アクションのイベントIDを返します
func (a Action) EventID() UUID {
	return a.eventID
}

// アクションの種類を返します
func (a Action) Kind() Kind {
	return a.kind
}

// アクションの実行する順番を返します
func (a Action) Order() Order {
	return a.order
}

// アクションの検証を行います
func (a Action) Validate() error {
	return nil
}

// 構造体からJSONに変換します
func (a Action) MarshalJSON() ([]byte, error) {
	data := struct {
		ID      UUID  `json:"id"`
		EventID UUID  `json:"event_id"`
		Kind    Kind  `json:"kind"`
		Order   Order `json:"order"`
	}{
		ID:      a.id,
		EventID: a.eventID,
		Kind:    a.kind,
		Order:   a.order,
	}

	return json.Marshal(data)
}

// JSONから構造体に変換します
func (a *Action) UnmarshalJSON(b []byte) error {
	data := struct {
		ID      UUID  `json:"id"`
		EventID UUID  `json:"event_id"`
		Kind    Kind  `json:"kind"`
		Order   Order `json:"order"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONのUnmarshalに失敗しました", err)
	}

	a.id = data.ID
	a.eventID = data.EventID
	a.kind = data.Kind
	a.order = data.Order

	if err := a.Validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}
