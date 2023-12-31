package action

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// 埋め込みのタイトルです
type Title struct {
	value string
}

// 埋め込みのタイトルを生成します
func NewTitle(value string) (Title, error) {
	t := Title{value: value}

	if err := t.validate(); err != nil {
		return Title{}, errors.NewError("埋め込みのタイトルが不正です", err)
	}

	return t, nil
}

// 埋め込みのタイトルを返します
func (t Title) String() string {
	return t.value
}

// 埋め込みのタイトルが存在しているか確認します
func (t Title) IsEmpty() bool {
	return t.value == ""
}

// 埋め込みのタイトルを検証します
func (t Title) validate() error {
	if len([]rune(t.value)) > 256 {
		return errors.NewError("送信埋め込みのタイトルの最大文字数を超えています")
	}

	return nil
}

// 埋め込みのタイトルをJSONに変換します
func (t Title) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: t.value,
	}

	return json.Marshal(data)
}

// 埋め込みのタイトルをJSONから変換します
func (t *Title) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	t.value = data.Value

	return nil
}
