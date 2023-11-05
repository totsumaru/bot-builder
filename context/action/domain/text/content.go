package text

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const ContentMaxLength = 1500

// 内容です
type Content struct {
	value string
}

// 内容を生成します
func NewContent(value string) (Content, error) {
	t := Content{value: value}

	if err := t.validate(); err != nil {
		return t, errors.NewError("検証に失敗しました", err)
	}

	return t, nil
}

// 内容を返します
func (t Content) String() string {
	return t.value
}

// 内容が存在しているか確認します
func (t Content) IsEmpty() bool {
	return t.value == ""
}

// 内容を検証します
func (t Content) validate() error {
	if len([]rune(t.value)) > ContentMaxLength {
		return errors.NewError("Contentの最大文字数を超えています")
	}

	return nil
}

// 内容をJSONに変換します
func (t Content) MarshalJSON() ([]byte, error) {
	data := struct {
		Content string `json:"content"`
	}{
		Content: t.value,
	}

	return json.Marshal(data)
}

// JSONから内容を復元します
func (t *Content) UnmarshalJSON(b []byte) error {
	data := struct {
		Content string `json:"content"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからContentの復元に失敗しました", err)
	}

	t.value = data.Content

	return nil
}
