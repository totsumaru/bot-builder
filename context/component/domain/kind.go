package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	ComponentKindButton = "BUTTON"
	ComponentKindImage  = "IMAGE"
)

// コンポーネントの種類です
type Kind struct {
	value string
}

// コンポーネントの種類を作成します
func NewKind(value string) (Kind, error) {
	k := Kind{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// コンポーネントの種類を返します
func (k Kind) String() string {
	return k.value
}

// コンポーネントの種類が存在しているか確認します
func (k Kind) IsEmpty() bool {
	return k.value == ""
}

// コンポーネントの種類を検証します
func (k Kind) validate() error {
	switch k.value {
	case ComponentKindButton:
	default:
		return errors.NewError("コンポーネントの種類が不正です")
	}

	return nil
}

// コンポーネントの種類をJSONに変換します
func (k Kind) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: k.value,
	}

	return json.Marshal(data)
}

// コンポーネントの種類をJSONから復元します
func (k *Kind) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	k.value = data.Value

	return nil
}
