package reply_text

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// テキストを返信するアクションです
type ReplyText struct {
	actionType  action.ActionType
	content     action.Content
	isEphemeral bool
	componentID []context.UUID // コンポーネントのID
}

// テキストを返信するアクションを生成します
func NewReplyText(
	content action.Content,
	isEphemeral bool,
	componentID []context.UUID,
) (ReplyText, error) {
	at, err := action.NewActionType(action.ActionTypeReplyText)
	if err != nil {
		return ReplyText{}, errors.NewError("アクションタイプを作成できません", err)
	}

	s := ReplyText{
		actionType:  at,
		content:     content,
		isEphemeral: isEphemeral,
		componentID: componentID,
	}

	if err = s.validate(); err != nil {
		return ReplyText{}, err
	}

	return s, nil
}

// アクションタイプを返します
func (s ReplyText) ActionType() action.ActionType {
	return s.actionType
}

// 送信する内容を返します
func (s ReplyText) Content() action.Content {
	return s.content
}

// エフェメラルかどうかを返します
func (s ReplyText) IsEphemeral() bool {
	return s.isEphemeral
}

// コンポーネントのIDを返します
func (s ReplyText) ComponentID() []context.UUID {
	return s.componentID
}

// 検証します
func (s ReplyText) validate() error {
	if len(s.componentID) > 5 {
		return errors.NewError("コンポーネントIDは5つまでです", nil)
	}

	return nil
}

// JSONに変換します
func (s ReplyText) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType  action.ActionType `json:"action_type"`
		Content     action.Content    `json:"content"`
		IsEphemeral bool              `json:"is_ephemeral"`
		ComponentID []context.UUID    `json:"component_id"`
	}{
		ActionType:  s.actionType,
		Content:     s.content,
		IsEphemeral: s.isEphemeral,
		ComponentID: s.componentID,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (s *ReplyText) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType  action.ActionType `json:"action_type"`
		Content     action.Content    `json:"content"`
		IsEphemeral bool              `json:"is_ephemeral"`
		ComponentID []context.UUID    `json:"component_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	s.actionType = data.ActionType
	s.content = data.Content
	s.isEphemeral = data.IsEphemeral
	s.componentID = data.ComponentID

	return nil
}
