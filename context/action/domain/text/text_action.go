package text

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/action/domain"
	"github.com/totsumaru/bot-builder/context/action/domain/components/button"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// テキストアクションです
type TextAction struct {
	domain.ActionCore
	content     Content
	button      []button.Button
	isResponse  bool
	isEphemeral bool
	channelID   domain.DiscordID // 返信の場合は必ず空になる
}

// テキストアクションを作成します
func NewTextAction(
	eventID domain.UUID,
	order domain.Order,
	content Content,
	button []button.Button,
	isResponse bool,
	isEphemeral bool,
	channelID domain.DiscordID,
) (TextAction, error) {
	k, err := domain.NewKind(domain.ActionKindText)
	if err != nil {
		return TextAction{}, errors.NewError("アクションの種類の生成に失敗しました", err)
	}

	act, err := domain.NewAction(eventID, k, order)
	if err != nil {
		return TextAction{}, errors.NewError("アクションの生成に失敗しました", err)
	}

	ta := TextAction{
		ActionCore:  act,
		content:     content,
		button:      button,
		isResponse:  isResponse,
		isEphemeral: isEphemeral,
		channelID:   channelID,
	}

	if err = ta.validate(); err != nil {
		return ta, errors.NewError("検証に失敗しました", err)
	}

	return ta, nil
}

// テキストアクションの内容を返します
func (ta TextAction) Content() Content {
	return ta.content
}

// テキストアクションのボタンを返します
func (ta TextAction) Button() []button.Button {
	return ta.button
}

// テキストアクションが返信か確認します
func (ta TextAction) IsResponse() bool {
	return ta.isResponse
}

// テキストアクションがエフェメラルか確認します
func (ta TextAction) IsEphemeral() bool {
	return ta.isEphemeral
}

// テキストアクションのチャンネルIDを返します
func (ta TextAction) ChannelID() domain.DiscordID {
	return ta.channelID
}

// テキストアクションを検証します
func (ta TextAction) validate() error {
	// 返信の場合は必ず空の値
	if ta.isResponse {
		if ta.channelID.String() != "" {
			return errors.NewError("返信の場合はチャンネルIDは空にしてください")
		}
	}

	return nil
}

// テキストアクションをJSONに変換します
func (ta TextAction) MarshalJSON() ([]byte, error) {
	data := struct {
		ID          domain.UUID      `json:"id"`
		EventID     domain.UUID      `json:"event_id"`
		Kind        domain.Kind      `json:"kind"`
		Order       domain.Order     `json:"order"`
		Content     Content          `json:"content"`
		Button      []button.Button  `json:"button"`
		IsResponse  bool             `json:"is_response"`
		IsEphemeral bool             `json:"is_ephemeral"`
		ChannelID   domain.DiscordID `json:"channel_id"`
	}{
		ID:          ta.ID(),
		EventID:     ta.EventID(),
		Kind:        ta.Kind(),
		Order:       ta.Order(),
		Content:     ta.content,
		Button:      ta.button,
		IsResponse:  ta.isResponse,
		IsEphemeral: ta.isEphemeral,
		ChannelID:   ta.channelID,
	}

	return json.Marshal(data)
}

// テキストアクションをJSONから復元します
func (ta *TextAction) UnmarshalJSON(b []byte) error {
	data := struct {
		ID          domain.UUID      `json:"id"`
		EventID     domain.UUID      `json:"event_id"`
		Kind        domain.Kind      `json:"kind"`
		Order       domain.Order     `json:"order"`
		Content     Content          `json:"content"`
		Button      []button.Button  `json:"button"`
		IsResponse  bool             `json:"is_response"`
		IsEphemeral bool             `json:"is_ephemeral"`
		ChannelID   domain.DiscordID `json:"channel_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからテキストアクションを復元できませんでした", err)
	}

	act, err := domain.NewAction(data.ID, data.Kind, data.Order)
	if err != nil {
		return errors.NewError("JSONからテキストアクションを復元できませんでした", err)
	}

	ta.ActionCore = act
	ta.content = data.Content
	ta.button = data.Button
	ta.isResponse = data.IsResponse
	ta.isEphemeral = data.IsEphemeral
	ta.channelID = data.ChannelID

	return nil
}
