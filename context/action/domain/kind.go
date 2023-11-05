package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// Actionの種類を表す定数です
type Kind struct {
	value string
}

const (
	ActionKindText  = "TEXT"  // 通常のテキストのActionです
	ActionKindEmbed = "EMBED" // EmbedのActionです
	// TODO:　他にもActionの種類を追加していく
)

// Actionの種類を作成します
func NewKind(value string) (Kind, error) {
	k := Kind{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// Actionの種類を返します
func (k Kind) String() string {
	return k.value
}

// Actionの種類が存在しているか確認します
func (k Kind) IsEmpty() bool {
	return k.value == ""
}

// Actionの種類を検証します
func (k Kind) validate() error {
	switch k.value {
	case ActionKindText:
	case ActionKindEmbed:
	default:
		return errors.NewError("Actionの種類が不正です")
	}

	return nil
}

// Actionの種類をJSONに変換します
func (k Kind) MarshalJSON() ([]byte, error) {
	data := struct {
		Kind string `json:"kind"`
	}{
		Kind: k.value,
	}

	return json.Marshal(data)
}

// Actionの種類をJSONから復元します
func (k *Kind) UnmarshalJSON(b []byte) error {
	data := struct {
		Kind string `json:"kind"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	k.value = data.Kind

	if err := k.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}
