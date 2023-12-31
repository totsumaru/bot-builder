package button

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const LabelMaxLength = 8

// ラベルです
type Label struct {
	value string
}

// ラベルを生成します
func NewLabel(value string) (Label, error) {
	l := Label{value: value}

	if err := l.validate(); err != nil {
		return l, errors.NewError("検証に失敗しました", err)
	}

	return l, nil
}

// ラベルを返します
func (l Label) String() string {
	return l.value
}

// ラベルが存在しているか確認します
func (l Label) IsEmpty() bool {
	return l.value == ""
}

// ラベルを検証します
func (l Label) validate() error {
	if l.IsEmpty() {
		return errors.NewError("Labelが空です")
	}

	if len([]rune(l.value)) > LabelMaxLength {
		return errors.NewError("ラベルの最大文字数を超えています")
	}

	return nil
}

// ラベルをJSONに変換します
func (l Label) MarshalJSON() ([]byte, error) {
	data := struct {
		Label string `json:"label"`
	}{
		Label: l.value,
	}

	return json.Marshal(data)
}

// JSONからラベルを復元します
func (l *Label) UnmarshalJSON(b []byte) error {
	data := struct {
		Label string `json:"label"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからLabelの復元に失敗しました", err)
	}

	l.value = data.Label

	return nil
}
