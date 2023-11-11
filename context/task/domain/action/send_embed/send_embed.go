package send_embed

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// Embedのテキストを送信するアクションです
type SendEmbed struct {
	actionType       action.ActionType
	channelID        context.DiscordID
	title            action.Title
	content          action.Content
	colorCode        action.ColorCode
	imageComponentID context.UUID
	displayAuthor    bool
	componentID      []context.UUID
}

// Embedのテキストを送信するアクションを生成します
func NewSendEmbed(
	channelID context.DiscordID,
	title action.Title,
	content action.Content,
	colorCode action.ColorCode,
	imageComponentID context.UUID,
	displayAuthor bool,
	componentID []context.UUID,
) (SendEmbed, error) {
	at, err := action.NewActionType(action.ActionTypeSendEmbed)
	if err != nil {
		return SendEmbed{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := SendEmbed{
		actionType:       at,
		channelID:        channelID,
		title:            title,
		content:          content,
		colorCode:        colorCode,
		imageComponentID: imageComponentID,
		displayAuthor:    displayAuthor,
		componentID:      componentID,
	}

	if err = s.validate(); err != nil {
		return SendEmbed{}, err
	}

	return s, nil
}

// アクションタイプを返します
func (s SendEmbed) ActionType() action.ActionType {
	return s.actionType
}

// チャンネルIDを返します
func (s SendEmbed) ChannelID() context.DiscordID {
	return s.channelID
}

// タイトルを返します
func (s SendEmbed) Title() action.Title {
	return s.title
}

// 内容を返します
func (s SendEmbed) Content() action.Content {
	return s.content
}

// カラーコードを返します
func (s SendEmbed) ColorCode() action.ColorCode {
	return s.colorCode
}

// 画像URLを返します
func (s SendEmbed) ImageComponentID() context.UUID {
	return s.imageComponentID
}

// Authorを表示するかどうかを返します
func (s SendEmbed) DisplayAuthor() bool {
	return s.displayAuthor
}

// コンポーネントのIDを返します
func (s SendEmbed) ComponentID() []context.UUID {
	return s.componentID
}

// 検証します
func (s SendEmbed) validate() error {
	return nil
}

// JSONに変換します
func (s SendEmbed) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType       action.ActionType `json:"action_type"`
		ChannelID        context.DiscordID `json:"channel_id"`
		Title            action.Title      `json:"title"`
		Content          action.Content    `json:"content"`
		ColorCode        action.ColorCode  `json:"color_code"`
		ImageComponentID context.UUID      `json:"image_component_id"`
		DisplayAuthor    bool              `json:"display_author"`
		ComponentID      []context.UUID    `json:"component_id"`
	}{
		ActionType:       s.actionType,
		ChannelID:        s.channelID,
		Title:            s.title,
		Content:          s.content,
		ColorCode:        s.colorCode,
		ImageComponentID: s.imageComponentID,
		DisplayAuthor:    s.displayAuthor,
		ComponentID:      s.componentID,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *SendEmbed) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType       action.ActionType `json:"action_type"`
		ChannelID        context.DiscordID `json:"channel_id"`
		Title            action.Title      `json:"title"`
		Content          action.Content    `json:"content"`
		ColorCode        action.ColorCode  `json:"color_code"`
		ImageComponentID context.UUID      `json:"image_component_id"`
		DisplayAuthor    bool              `json:"display_author"`
		ComponentID      []context.UUID    `json:"component_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.actionType = data.ActionType
	s.channelID = data.ChannelID
	s.title = data.Title
	s.content = data.Content
	s.colorCode = data.ColorCode
	s.imageComponentID = data.ImageComponentID
	s.displayAuthor = data.DisplayAuthor
	s.componentID = data.ComponentID

	return nil
}
