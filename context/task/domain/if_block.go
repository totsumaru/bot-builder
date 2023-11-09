package domain

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/task/domain/condition"
)

// ifブロックです
//
// trueAction/falseActionにはIfBlock/Actionのいずれかが入ります。
type IfBlock struct {
	condition   condition.Condition // 条件
	trueAction  []any               // 条件がtrueの場合のアクション
	falseAction []any               // 条件がfalseの場合のアクション
}

// ifブロックを生成します
func NewIfBlock(
	condition condition.Condition,
	trueAction []any,
	falseAction []any,
) (IfBlock, error) {
	ifBlock := IfBlock{
		condition:   condition,
		trueAction:  trueAction,
		falseAction: falseAction,
	}

	if err := ifBlock.validate(); err != nil {
		return IfBlock{}, err
	}

	return ifBlock, nil
}

// 条件を返します
func (i IfBlock) Condition() condition.Condition {
	return i.condition
}

// 条件がtrueの場合のアクションを返します
func (i IfBlock) TrueAction() []any {
	return i.trueAction
}

// 条件がfalseの場合のアクションを返します
func (i IfBlock) FalseAction() []any {
	return i.falseAction
}

// 検証します
func (i IfBlock) validate() error {
	return nil
}

// JSONに変換します
func (i IfBlock) MarshalJSON() ([]byte, error) {
	data := struct {
		Condition   condition.Condition `json:"condition"`
		TrueAction  []any               `json:"true_action"`
		FalseAction []any               `json:"false_action"`
	}{
		Condition:   i.condition,
		TrueAction:  i.trueAction,
		FalseAction: i.falseAction,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (i *IfBlock) UnmarshalJSON(b []byte) error {
	data := struct {
		Condition   condition.Condition `json:"condition"`
		TrueAction  []any               `json:"true_action"`
		FalseAction []any               `json:"false_action"`
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

	i.condition = ifBlock.condition
	i.trueAction = ifBlock.trueAction
	i.falseAction = ifBlock.falseAction

	return nil
}
