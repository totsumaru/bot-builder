package button

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// ボタンコンポーネントです
type Button struct {
	domain.ComponentCore
	label Label
	style Style
	url   context.URL
}

// ボタンを作成します
func NewButton(
	core domain.ComponentCore,
	label Label,
	style Style,
	url context.URL,
) (Button, error) {
	b := Button{
		ComponentCore: core,
		label:         label,
		style:         style,
		url:           url,
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

// ボタンのURLを返します
func (b Button) URL() context.URL {
	return b.url
}

// 検証します
func (b Button) validate() error {
	// Linkの時はURLが必須
	switch b.style.String() {
	case ButtonStyleLink:
		if b.url.IsEmpty() {
			return errors.NewError("URLが空です")
		}
	}

	return nil
}

// MarshalJSON は Button 構造体を JSON に変換します。
func (b Button) MarshalJSON() ([]byte, error) {
	// ComponentCore のフィールドを直接含める
	bb, err := json.Marshal(struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		Kind          domain.Kind       `json:"kind"`
		Label         Label             `json:"label"`
		Style         Style             `json:"style"`
		URL           context.URL       `json:"url"`
	}{
		ID:            b.ID(),
		ServerID:      b.ServerID(),
		ApplicationID: b.ApplicationID(),
		Kind:          b.Kind(),
		Label:         b.label,
		Style:         b.style,
		URL:           b.url,
	})

	if err != nil {
		return nil, errors.NewError("JSONに変換できませんでした", err)
	}

	return bb, nil
}

// UnmarshalJSON は JSON から Button 構造体を復元します。
func (b *Button) UnmarshalJSON(bytes []byte) error {
	// ComponentCore を含めた一時的な構造体を定義
	data := struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		Kind          domain.Kind       `json:"kind"`
		Label         Label             `json:"label"`
		Style         Style             `json:"style"`
		URL           context.URL       `json:"url"`
	}{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.NewError("JSONからボタンの復元に失敗しました", err)
	}

	id, err := context.RestoreUUID(data.ID.String())
	if err != nil {
		return errors.NewError("IDを復元できません", err)
	}

	serverID, err := context.NewDiscordID(data.ServerID.String())
	if err != nil {
		return errors.NewError("DiscordIDを作成できません", err)
	}

	appID, err := context.RestoreUUID(data.ApplicationID.String())
	if err != nil {
		return errors.NewError("ApplicationIDを復元できません", err)
	}

	kind, err := domain.NewKind(data.Kind.String())
	if err != nil {
		return errors.NewError("コンポーネントの種類を作成できません", err)
	}

	b.ComponentCore, err = domain.NewComponentCore(id, serverID, appID, kind)
	if err != nil {
		return errors.NewError("コンポーネントの共通の構造体を作成できません", err)
	}

	b.label = data.Label
	b.style = data.Style
	b.url = data.URL

	return nil
}
