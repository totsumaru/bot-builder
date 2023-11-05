package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	EventKindMessageCreate = "MESSAGE" // メッセージが送信された時のイベントです
	EventKindButton        = "BUTTON"  // ボタンが押された時のイベントです
)

// Eventの種類です
type Kind struct {
	value string
}

// イベントの種類を作成します
func NewKind(value string) (Kind, error) {
	k := Kind{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// イベントの種類を返します
func (k Kind) String() string {
	return k.value
}

// イベントの種類が存在しているか確認します
func (k Kind) IsEmpty() bool {
	return k.value == ""
}

// イベントの種類を検証します
func (k Kind) validate() error {
	switch k.value {
	case EventKindMessageCreate:
	case EventKindButton:
	default:
		return errors.NewError("イベントの種類が不正です")
	}

	return nil
}

// イベントの種類をJSONに変換します
func (k Kind) MarshalJSON() ([]byte, error) {
	data := struct {
		Kind string `json:"kind"`
	}{
		Kind: k.value,
	}

	return json.Marshal(data)
}

// JSONからイベントの種類を復元します
func (k *Kind) UnmarshalJSON(b []byte) error {
	data := struct {
		Kind string `json:"kind"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからイベントの種類の復元に失敗しました", err)
	}

	k.value = data.Kind

	return nil
}
