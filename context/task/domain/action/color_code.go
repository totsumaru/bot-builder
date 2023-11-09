package action

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	ColorBlue     = 0x0099ff
	ColorRed      = 0xff0000
	ColorOrange   = 0xffa500
	ColorGreen    = 0x3cb371
	ColorPink     = 0xff69b4
	ColorBlack    = 0x000000
	ColorYellow   = 0xffd700
	ColorCyan     = 0x00ffff
	ColorDarkGray = 0x2c2c34
)

// カラーコードです
type ColorCode struct {
	value int
}

// カラーコードを生成します
func NewColorCode(value int) (ColorCode, error) {
	c := ColorCode{value: value}

	if err := c.validate(); err != nil {
		return ColorCode{}, errors.NewError("カラーコードが不正です", err)
	}

	return c, nil
}

// カラーコードを返します
func (c ColorCode) String() int {
	return c.value
}

// カラーコードが存在しているか確認します
func (c ColorCode) IsZero() bool {
	return c.value == 0
}

// カラーコードを検証します
func (c ColorCode) validate() error {
	return nil
}

// カラーコードをJSONに変換します
func (c ColorCode) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: c.value,
	}

	return json.Marshal(data)
}

// カラーコードをJSONから変換します
func (c *ColorCode) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	c.value = data.Value

	return nil
}
