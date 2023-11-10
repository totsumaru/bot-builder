package reply_embed

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// Embedのテキストを返信するアクションです
type ReplyEmbed struct {
	actionType    action.ActionType
	title         action.Title
	content       action.Content
	colorCode     action.ColorCode
	imageURL      domain.URL
	displayAuthor bool
	isEphemeral   bool
}

// Embedのテキストを返信するアクションを生成します
func NewReplyEmbed(
	title action.Title,
	content action.Content,
	colorCode action.ColorCode,
	imageURL domain.URL,
	displayAuthor bool,
	isEphemeral bool,
) (ReplyEmbed, error) {
	at, err := action.NewActionType(action.ActionTypeReplyEmbed)
	if err != nil {
		return ReplyEmbed{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := ReplyEmbed{
		actionType:    at,
		title:         title,
		content:       content,
		colorCode:     colorCode,
		imageURL:      imageURL,
		displayAuthor: displayAuthor,
		isEphemeral:   isEphemeral,
	}

	if err = s.validate(); err != nil {
		return ReplyEmbed{}, err
	}

	return s, nil
}

// アクションタイプを返します
func (s ReplyEmbed) ActionType() action.ActionType {
	return s.actionType
}

// タイトルを返します
func (s ReplyEmbed) Title() action.Title {
	return s.title
}

// 内容を返します
func (s ReplyEmbed) Content() action.Content {
	return s.content
}

// カラーコードを返します
func (s ReplyEmbed) ColorCode() action.ColorCode {
	return s.colorCode
}

// 画像URLを返します
func (s ReplyEmbed) ImageURL() domain.URL {
	return s.imageURL
}

// Authorを表示するかどうかを返します
func (s ReplyEmbed) DisplayAuthor() bool {
	return s.displayAuthor
}

// エフェメラルかどうかを返します
func (s ReplyEmbed) IsEphemeral() bool {
	return s.isEphemeral
}

// 検証します
func (s ReplyEmbed) validate() error {
	return nil
}

// JSONに変換します
func (s ReplyEmbed) MarshalJSON() ([]byte, error) {
	data := struct {
		Title         action.Title     `json:"title"`
		Content       action.Content   `json:"content"`
		ColorCode     action.ColorCode `json:"color_code"`
		ImageURL      domain.URL       `json:"image_url"`
		DisplayAuthor bool             `json:"display_author"`
		IsEphemeral   bool             `json:"is_ephemeral"`
	}{
		Title:         s.title,
		Content:       s.content,
		ColorCode:     s.colorCode,
		ImageURL:      s.imageURL,
		DisplayAuthor: s.displayAuthor,
		IsEphemeral:   s.isEphemeral,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *ReplyEmbed) UnmarshalJSON(b []byte) error {
	data := struct {
		Title         action.Title     `json:"title"`
		Content       action.Content   `json:"content"`
		ColorCode     action.ColorCode `json:"color_code"`
		ImageURL      domain.URL       `json:"image_url"`
		DisplayAuthor bool             `json:"display_author"`
		IsEphemeral   bool             `json:"is_ephemeral"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.title = data.Title
	s.content = data.Content
	s.colorCode = data.ColorCode
	s.imageURL = data.ImageURL
	s.displayAuthor = data.DisplayAuthor
	s.isEphemeral = data.IsEphemeral

	return nil
}
