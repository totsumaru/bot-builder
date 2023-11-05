package button

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const URLMaxLength = 200

// URLです
type URL struct {
	value string
}

// URLを生成します
func NewURL(value string) (URL, error) {
	u := URL{value: value}

	if err := u.validate(); err != nil {
		return u, errors.NewError("検証に失敗しました", err)
	}

	return u, nil
}

// URLを返します
func (u URL) String() string {
	return u.value
}

// URLが存在しているか確認します
func (u URL) IsEmpty() bool {
	return u.value == ""
}

// URLを検証します
func (u URL) validate() error {
	if len([]rune(u.value)) > URLMaxLength {
		return errors.NewError("URLの最大文字数を超えています")
	}

	return nil
}

// URLをJSONに変換します
func (u URL) MarshalJSON() ([]byte, error) {
	data := struct {
		URL string `json:"url"`
	}{
		URL: u.value,
	}

	return json.Marshal(data)
}

// JSONからURLを復元します
func (u *URL) UnmarshalJSON(b []byte) error {
	data := struct {
		URL string `json:"url"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからURLの復元に失敗しました", err)
	}

	u.value = data.URL

	return nil
}
