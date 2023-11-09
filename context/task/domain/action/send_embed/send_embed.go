package send_embed

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// Embedのテキストを送信するアクションです
type SendEmbed struct {
	actionType    action.ActionType
	channelID     domain.DiscordID
	title         action.Title
	content       action.Content
	colorCode     action.ColorCode
	imageURL      domain.URL
	displayAuthor bool
}

// Embedのテキストを送信するアクションを生成します
func NewSendEmbed(
	channelID domain.DiscordID,
	title action.Title,
	content action.Content,
	colorCode action.ColorCode,
	imageURL domain.URL,
	displayAuthor bool,
) (SendEmbed, error) {
	at, err := action.NewActionType(action.ActionTypeSendEmbed)
	if err != nil {
		return SendEmbed{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := SendEmbed{
		actionType:    at,
		channelID:     channelID,
		title:         title,
		content:       content,
		colorCode:     colorCode,
		imageURL:      imageURL,
		displayAuthor: displayAuthor,
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
func (s SendEmbed) ChannelID() domain.DiscordID {
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
func (s SendEmbed) ImageURL() domain.URL {
	return s.imageURL
}

// Authorを表示するかどうかを返します
func (s SendEmbed) DisplayAuthor() bool {
	return s.displayAuthor
}

// 検証します
func (s SendEmbed) validate() error {
	return nil
}

// JSONに変換します
func (s SendEmbed) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType    action.ActionType `json:"action_type"`
		ChannelID     domain.DiscordID  `json:"channel_id"`
		Title         action.Title      `json:"title"`
		Content       action.Content    `json:"content"`
		ColorCode     action.ColorCode  `json:"color_code"`
		ImageURL      domain.URL        `json:"image_url"`
		DisplayAuthor bool              `json:"display_author"`
	}{
		ActionType:    s.actionType,
		ChannelID:     s.channelID,
		Title:         s.title,
		Content:       s.content,
		ColorCode:     s.colorCode,
		ImageURL:      s.imageURL,
		DisplayAuthor: s.displayAuthor,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *SendEmbed) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType    action.ActionType `json:"action_type"`
		ChannelID     domain.DiscordID  `json:"channel_id"`
		Title         action.Title      `json:"title"`
		Content       action.Content    `json:"content"`
		ColorCode     action.ColorCode  `json:"color_code"`
		ImageURL      domain.URL        `json:"image_url"`
		DisplayAuthor bool              `json:"display_author"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.actionType = data.ActionType
	s.channelID = data.ChannelID
	s.title = data.Title
	s.content = data.Content
	s.colorCode = data.ColorCode
	s.imageURL = data.ImageURL
	s.displayAuthor = data.DisplayAuthor

	return nil
}
