package send_text

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// テキストを送信するアクションです
type SendText struct {
	actionType  action.ActionType
	channelID   context.DiscordID
	content     action.Content
	componentID []context.UUID
}

// テキストを送信するアクションを生成します
func NewSendText(
	channelID context.DiscordID,
	content action.Content,
	componentID []context.UUID,
) (SendText, error) {
	at, err := action.NewActionType(action.ActionTypeSendText)
	if err != nil {
		return SendText{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := SendText{
		actionType:  at,
		channelID:   channelID,
		content:     content,
		componentID: componentID,
	}

	if err = s.validate(); err != nil {
		return SendText{}, err
	}

	return s, nil
}

// アクションタイプを返します
func (s SendText) ActionType() action.ActionType {
	return s.actionType
}

// チャンネルIDを返します
func (s SendText) ChannelID() context.DiscordID {
	return s.channelID
}

// 送信する内容を返します
func (s SendText) Content() action.Content {
	return s.content
}

// コンポーネントのIDを返します
func (s SendText) ComponentID() []context.UUID {
	return s.componentID
}

// 検証します
func (s SendText) validate() error {
	if len(s.componentID) > 5 {
		return errors.NewError("コンポーネントIDは5つまでです", nil)
	}

	return nil
}

// JSONに変換します
func (s SendText) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType  action.ActionType `json:"action_type"`
		ChannelID   context.DiscordID `json:"channel_id"`
		Content     action.Content    `json:"content"`
		ComponentID []context.UUID    `json:"component_id"`
	}{
		ActionType:  s.actionType,
		ChannelID:   s.channelID,
		Content:     s.content,
		ComponentID: s.componentID,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *SendText) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType  action.ActionType `json:"action_type"`
		ChannelID   context.DiscordID `json:"channel_id"`
		Content     action.Content    `json:"content"`
		ComponentID []context.UUID    `json:"component_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.actionType = data.ActionType
	s.channelID = data.ChannelID
	s.content = data.Content
	s.componentID = data.ComponentID

	return nil
}
