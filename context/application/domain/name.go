package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// アプリケーションの名前です
type Name struct {
	value string
}

// アプリケーションの名前を生成します
func NewName(value string) (Name, error) {
	n := Name{value: value}

	if err := n.validate(); err != nil {
		return n, err
	}

	return n, nil
}

// アプリケーションの名前を返します
func (n Name) String() string {
	return n.value
}

// アプリケーションの名前が存在しているか確認します
func (n Name) IsEmpty() bool {
	return n.value == ""
}

// アプリケーションの名前を検証します
func (n Name) validate() error {
	if n.value == "" {
		return errors.NewError("アプリケーションの名前が空です")
	}

	return nil
}

// アプリケーションの名前をJSONに変換します
func (n Name) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: n.value,
	}

	return json.Marshal(data)
}

// JSONからアプリケーションの名前を復元します
func (n *Name) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからアプリケーションの名前を復元できません", err)
	}

	n.value = data.Value

	return nil
}
