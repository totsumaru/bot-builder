package action

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/lib/errors"
)

const (
	ActionTypeSendText   = "SEND_TEXT"
	ActionTypeSendEmbed  = "SEND_EMBED"
	ActionTypeReplyText  = "REPLY_TEXT"
	ActionTypeReplyEmbed = "REPLY_EMBED"
	ActionTypeIfBlock    = "IF_BLOCK"
)

// アクションタイプです
type ActionType struct {
	value string
}

// アクションタイプを生成します
func NewActionType(value string) (ActionType, error) {
	a := ActionType{
		value: value,
	}

	if err := a.validate(); err != nil {
		return ActionType{}, err
	}

	return a, nil
}

// アクションタイプを返します
func (a ActionType) String() string {
	return a.value
}

// 検証します
func (a ActionType) validate() error {
	switch a.value {
	case ActionTypeSendText,
		ActionTypeSendEmbed,
		ActionTypeReplyText,
		ActionTypeReplyEmbed,
		ActionTypeIfBlock:
	default:
		return errors.NewError("アクションタイプが不正です")
	}

	return nil
}

// JSONに変換します
func (a ActionType) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: a.value,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (a *ActionType) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	a.value = data.Value

	return nil
}
