package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// ifブロックです
//
// trueAction/falseActionにはIfBlock/Actionのいずれかが入ります。
// アクションの種類(IfBlockですが、trueAction/falseActionの中に入る可能性があるため、入れています）。
type IfBlock struct {
	actionType  action.ActionType
	condition   condition.Condition // 条件
	trueAction  []action.Action     // 条件がtrueの場合のアクション
	falseAction []action.Action     // 条件がfalseの場合のアクション
}

// ifブロックを生成します
func NewIfBlock(
	condition condition.Condition,
	trueAction []action.Action,
	falseAction []action.Action,
) (IfBlock, error) {
	at, err := action.NewActionType(action.ActionTypeIfBlock)
	if err != nil {
		return IfBlock{}, errors.NewError("アクションタイプを作成できません", err)
	}

	ifBlock := IfBlock{
		actionType:  at,
		condition:   condition,
		trueAction:  trueAction,
		falseAction: falseAction,
	}

	if err = ifBlock.validate(); err != nil {
		return IfBlock{}, err
	}

	return ifBlock, nil
}

// アクションの種類を返します
func (i IfBlock) ActionType() action.ActionType {
	return i.actionType
}

// 条件を返します
func (i IfBlock) Condition() condition.Condition {
	return i.condition
}

// 条件がtrueの場合のアクションを返します
func (i IfBlock) TrueAction() []action.Action {
	return i.trueAction
}

// 条件がfalseの場合のアクションを返します
func (i IfBlock) FalseAction() []action.Action {
	return i.falseAction
}

// 検証します
func (i IfBlock) validate() error {
	return nil
}

// JSONに変換します
func (i IfBlock) MarshalJSON() ([]byte, error) {
	data := struct {
		ActionType  action.ActionType   `json:"action_type"`
		Condition   condition.Condition `json:"condition"`
		TrueAction  []action.Action     `json:"true_action"`
		FalseAction []action.Action     `json:"false_action"`
	}{
		ActionType:  i.actionType,
		Condition:   i.condition,
		TrueAction:  i.trueAction,
		FalseAction: i.falseAction,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (i *IfBlock) UnmarshalJSON(b []byte) error {
	data := struct {
		ActionType  action.ActionType   `json:"action_type"`
		Condition   condition.Condition `json:"condition"`
		TrueAction  []action.Action     `json:"true_action"`
		FalseAction []action.Action     `json:"false_action"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	ifBlock, err := NewIfBlock(
		data.Condition,
		data.TrueAction,
		data.FalseAction,
	)
	if err != nil {
		return err
	}

	i.actionType = ifBlock.actionType
	i.condition = ifBlock.condition
	i.trueAction = ifBlock.trueAction
	i.falseAction = ifBlock.falseAction

	return nil
}
