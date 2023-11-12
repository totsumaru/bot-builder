package reply_embed

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// Embedのテキストを返信するアクションです
type ReplyEmbed struct {
	actionType       action.ActionType
	title            action.Title
	content          action.Content
	colorCode        action.ColorCode
	imageComponentID context.UUID
	displayAuthor    bool
	isEphemeral      bool
	componentID      []context.UUID
}

// Embedのテキストを返信するアクションを生成します
func NewReplyEmbed(
	title action.Title,
	content action.Content,
	colorCode action.ColorCode,
	imageComponentID context.UUID,
	displayAuthor bool,
	isEphemeral bool,
	componentID []context.UUID,
) (ReplyEmbed, error) {
	at, err := action.NewActionType(action.ActionTypeReplyEmbed)
	if err != nil {
		return ReplyEmbed{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := ReplyEmbed{
		actionType:       at,
		title:            title,
		content:          content,
		colorCode:        colorCode,
		imageComponentID: imageComponentID,
		displayAuthor:    displayAuthor,
		isEphemeral:      isEphemeral,
		componentID:      componentID,
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
func (s ReplyEmbed) ImageComponentID() context.UUID {
	return s.imageComponentID
}

// Authorを表示するかどうかを返します
func (s ReplyEmbed) DisplayAuthor() bool {
	return s.displayAuthor
}

// エフェメラルかどうかを返します
func (s ReplyEmbed) IsEphemeral() bool {
	return s.isEphemeral
}

// コンポーネントのIDを返します
func (s ReplyEmbed) ComponentID() []context.UUID {
	return s.componentID
}

// 検証します
func (s ReplyEmbed) validate() error {
	if len(s.componentID) > 5 {
		return errors.NewError("コンポーネントIDは5つまでです", nil)
	}

	return nil
}

// JSONに変換します
func (s ReplyEmbed) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType       action.ActionType `json:"action_type"`
		Title            action.Title      `json:"title"`
		Content          action.Content    `json:"content"`
		ColorCode        action.ColorCode  `json:"color_code"`
		ImageComponentID context.UUID      `json:"image_component_id"`
		DisplayAuthor    bool              `json:"display_author"`
		IsEphemeral      bool              `json:"is_ephemeral"`
		ComponentID      []context.UUID    `json:"component_id"`
	}{
		ActionType:       s.actionType,
		Title:            s.title,
		Content:          s.content,
		ColorCode:        s.colorCode,
		ImageComponentID: s.imageComponentID,
		DisplayAuthor:    s.displayAuthor,
		IsEphemeral:      s.isEphemeral,
		ComponentID:      s.componentID,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *ReplyEmbed) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType       action.ActionType `json:"action_type"`
		Title            action.Title      `json:"title"`
		Content          action.Content    `json:"content"`
		ColorCode        action.ColorCode  `json:"color_code"`
		ImageComponentID context.UUID      `json:"image_component_id"`
		DisplayAuthor    bool              `json:"display_author"`
		IsEphemeral      bool              `json:"is_ephemeral"`
		ComponentID      []context.UUID    `json:"component_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.actionType = data.ActionType
	s.title = data.Title
	s.content = data.Content
	s.colorCode = data.ColorCode
	s.imageComponentID = data.ImageComponentID
	s.displayAuthor = data.DisplayAuthor
	s.isEphemeral = data.IsEphemeral
	s.componentID = data.ComponentID

	return nil
}
