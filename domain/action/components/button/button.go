package button

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// ボタンです
type Button struct {
	label   Label
	style   Style
	eventID domain.UUID
	url     URL
}

// ボタンを作成します
func NewButton(label Label, style Style, eventID domain.UUID, url URL) (Button, error) {
	b := Button{
		label:   label,
		style:   style,
		eventID: eventID,
		url:     url,
	}

	if err := b.validate(); err != nil {
		return b, errors.NewError("検証に失敗しました", err)
	}

	return b, nil
}

// ボタンのラベルを返します
func (b Button) Label() Label {
	return b.label
}

// ボタンのスタイルを返します
func (b Button) Style() Style {
	return b.style
}

// ボタンのイベントIDを返します
func (b Button) EventID() domain.UUID {
	return b.eventID
}

// ボタンのURLを返します
func (b Button) URL() URL {
	return b.url
}

// 検証します
func (b Button) validate() error {
	switch b.style.String() {
	case ButtonStyleLink:
		if b.url.IsEmpty() {
			return errors.NewError("URLが空です")
		}
	}

	return nil
}

// ボタンをJSONに変換します
func (b Button) MarshalJSON() ([]byte, error) {
	data := struct {
		Label   Label       `json:"label"`
		Style   Style       `json:"style"`
		EventID domain.UUID `json:"event_id"`
		URL     URL         `json:"url"`
	}{
		Label:   b.label,
		Style:   b.style,
		EventID: b.eventID,
		URL:     b.url,
	}

	return json.Marshal(data)
}

// JSONからボタンを復元します
func (b *Button) UnmarshalJSON(bytes []byte) error {
	data := struct {
		Label   Label       `json:"label"`
		Style   Style       `json:"style"`
		EventID domain.UUID `json:"event_id"`
		URL     URL         `json:"url"`
	}{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.NewError("JSONからボタンの復元に失敗しました", err)
	}

	b.label = data.Label
	b.style = data.Style
	b.eventID = data.EventID
	b.url = data.URL

	return nil
}
