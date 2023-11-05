package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// 実行する順番です
type Order struct {
	value int
}

// 実行する順番を生成します
func NewOrder(value int) (Order, error) {
	o := Order{
		value: value,
	}

	if err := o.validate(); err != nil {
		return o, errors.NewError("検証に失敗しました", err)
	}

	return o, nil
}

// 実行する順番を返します
func (o Order) Int() int {
	return o.value
}

// 実行する順番が存在しているか確認します
func (o Order) IsZero() bool {
	return o.value == 0
}

// 実行する順番を検証します
func (o Order) validate() error {
	return nil
}

// 実行する順番をJSONに変換します
func (o Order) MarshalJSON() ([]byte, error) {
	data := struct {
		Order int `json:"order"`
	}{
		Order: o.value,
	}

	return json.Marshal(data)
}

// JSONから実行する順番を復元します
func (o *Order) UnmarshalJSON(b []byte) error {
	data := struct {
		Order int `json:"order"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからOrderの復元に失敗しました", err)
	}

	o.value = data.Order

	return nil
}
