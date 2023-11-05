package message

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const KeywordMaxLength = 30

// キーワードです
type Keyword struct {
	value string
}

// キーワードを生成します
func NewKeyword(value string) (Keyword, error) {
	d := Keyword{value: value}

	if err := d.validate(); err != nil {
		return d, errors.NewError("検証に失敗しました", err)
	}

	return d, nil
}

// キーワードを返します
func (k Keyword) String() string {
	return k.value
}

// キーワードが存在しているか確認します
func (k Keyword) IsEmpty() bool {
	return k.value == ""
}

// キーワードを検証します
func (k Keyword) validate() error {
	if k.IsEmpty() {
		return errors.NewError("Keywordが空です")
	}

	if len([]rune(k.value)) > KeywordMaxLength {
		return errors.NewError("キーワードの最大文字数を超えています")
	}

	return nil
}

// キーワードをJSONに変換します
func (k Keyword) MarshalJSON() ([]byte, error) {
	data := struct {
		Keyword string `json:"keyword"`
	}{
		Keyword: k.value,
	}

	return json.Marshal(data)
}

// JSONからキーワードを復元します
func (k *Keyword) UnmarshalJSON(b []byte) error {
	data := struct {
		Keyword string `json:"keyword"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからKeywordの復元に失敗しました", err)
	}

	k.value = data.Keyword

	return nil
}
