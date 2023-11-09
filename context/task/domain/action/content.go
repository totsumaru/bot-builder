package action

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

// 内容です
type Content struct {
	value string
}

// 内容を生成します
func NewContent(value string) (Content, error) {
	c := Content{value: value}

	if err := c.validate(); err != nil {
		return Content{}, errors.NewError("内容が不正です", err)
	}

	return c, nil
}

// 内容を返します
func (c Content) String() string {
	return c.value
}

// 内容が存在しているか確認します
func (c Content) IsEmpty() bool {
	return c.value == ""
}

// 内容を検証します
func (c Content) validate() error {
	if len([]rune(c.value)) > 1500 {
		return errors.NewError("送信内容の最大文字数を超えています")
	}

	return nil
}

// 内容をJSONに変換します
func (c Content) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: c.value,
	}

	return json.Marshal(data)
}

// 内容をJSONから変換します
func (c *Content) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	c.value = data.Value

	return nil
}
