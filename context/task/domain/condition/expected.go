package condition

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// 期待する値です
type Expected struct {
	value string
}

// 期待する値を作成します
func NewExpected(value string) (Expected, error) {
	k := Expected{value: value}

	if err := k.validate(); err != nil {
		return k, errors.NewError("検証に失敗しました", err)
	}

	return k, nil
}

// 期待する値を返します
func (e Expected) String() string {
	return e.value
}

// 期待する値が存在しているか確認します
func (e Expected) IsEmpty() bool {
	return e.value == ""
}

// 期待する値を検証します
//
// 様々な値が入る可能性があるため、validationはは緩めに設定します。
func (e Expected) validate() error {
	if e.IsEmpty() {
		return errors.NewError("期待する値が存在しません", nil)
	}

	return nil
}

// 期待する値をJSONに変換します
func (e Expected) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: e.value,
	}

	return json.Marshal(data)
}

// 期待する値をJSONから復元します
func (e *Expected) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	e.value = data.Value

	if err := e.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}
