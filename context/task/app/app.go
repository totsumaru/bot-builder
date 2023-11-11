package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_text"
	"github.com/totsumaru/bot-builder/context/task/domain/action/send_embed"
	"github.com/totsumaru/bot-builder/context/task/domain/action/send_text"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// ifブロックのリクエストです
type IfBlockReq struct {
	Condition struct {
		Kind     string
		Expected string
	}
	TrueAction  []any
	FalseAction []any
}

// ==============================================
// Actionのリクエスト
// ==============================================

// テキストを送信するアクションのリクエストです
type SendTextActionReq struct {
	ChannelID   string
	Content     string
	ComponentID []string
}

// テキストを返信するアクションのリクエストです
type ReplyTextActionReq struct {
	Content     string
	IsEphemeral bool
	ComponentID []string
}

// Embedを送信するアクションのリクエストです
type SendEmbedActionReq struct {
	ChannelID        string
	Title            string
	Content          string
	ColorCode        int
	ImageComponentID string
	DisplayAuthor    bool
}

// Embedを返信するアクションのリクエストです
type ReplyEmbedActionReq struct {
	Title            string
	Content          string
	ColorCode        int
	ImageComponentID string
	DisplayAuthor    bool
	IsEphemeral      bool
}

// ==============================================
// 共通処理
// ==============================================

// ifBlockを作成します
//
// Actionは再帰を使用します。
func CreateIfBlock(req IfBlockReq) (domain.IfBlock, error) {
	condKind, err := condition.NewKind(req.Condition.Kind)
	if err != nil {
		return domain.IfBlock{}, errors.NewError("条件の種類を作成できません", err)
	}

	expected, err := condition.NewExpected(req.Condition.Expected)
	if err != nil {
		return domain.IfBlock{}, errors.NewError("条件の期待値を作成できません", err)
	}

	cond, err := condition.NewCondition(condKind, expected)
	if err != nil {
		return domain.IfBlock{}, errors.NewError("条件を作成できません", err)
	}

	// trueActionを作成
	var trueActions []action.Action
	for _, reqAct := range req.TrueAction {
		act, err := CreateActionFromReq(reqAct)
		if err != nil {
			return domain.IfBlock{}, err
		}
		trueActions = append(trueActions, act)
	}

	// falseActionを作成
	var falseActions []action.Action
	for _, reqAct := range req.FalseAction {
		act, err := CreateActionFromReq(reqAct)
		if err != nil {
			return domain.IfBlock{}, err
		}
		falseActions = append(falseActions, act)
	}

	ifBlock, err := domain.NewIfBlock(cond, trueActions, falseActions)
	if err != nil {
		return domain.IfBlock{}, errors.NewError("IfBlockを作成できません", err)
	}

	return ifBlock, nil
}

// リクエストからActionを作成します
func CreateActionFromReq(req any) (action.Action, error) {
	switch reqTyped := req.(type) {
	case SendTextActionReq:
		chID, err := context.NewDiscordID(reqTyped.ChannelID)
		if err != nil {
			return nil, errors.NewError("DiscordIDを作成できません", err)
		}

		c, err := action.NewContent(reqTyped.Content)
		if err != nil {
			return nil, errors.NewError("Contentを作成できません", err)
		}

		componentID := make([]context.UUID, 0)
		for _, id := range reqTyped.ComponentID {
			cpID, err := context.RestoreUUID(id)
			if err != nil {
				return nil, errors.NewError("UUIDを作成できません", err)
			}
			componentID = append(componentID, cpID)
		}

		sendText, err := send_text.NewSendText(chID, c, componentID)
		if err != nil {
			return nil, errors.NewError("テキストアクションを作成できません", err)
		}

		return sendText, nil
	case ReplyTextActionReq:
		c, err := action.NewContent(reqTyped.Content)
		if err != nil {
			return nil, errors.NewError("Contentを作成できません", err)
		}

		componentID := make([]context.UUID, 0)
		for _, id := range reqTyped.ComponentID {
			cpID, err := context.RestoreUUID(id)
			if err != nil {
				return nil, errors.NewError("UUIDを作成できません", err)
			}
			componentID = append(componentID, cpID)
		}

		replyText, err := reply_text.NewReplyText(c, reqTyped.IsEphemeral, componentID)
		if err != nil {
			return nil, errors.NewError("テキストアクションを作成できません", err)
		}

		return replyText, nil
	case SendEmbedActionReq:
		chID, err := context.NewDiscordID(reqTyped.ChannelID)
		if err != nil {
			return nil, errors.NewError("DiscordIDを作成できません", err)
		}

		title, err := action.NewTitle(reqTyped.Title)
		if err != nil {
			return nil, errors.NewError("Titleを作成できません", err)
		}

		content, err := action.NewContent(reqTyped.Content)
		if err != nil {
			return nil, errors.NewError("Contentを作成できません", err)
		}

		colorCode, err := action.NewColorCode(reqTyped.ColorCode)
		if err != nil {
			return nil, errors.NewError("ColorCodeを作成できません", err)
		}

		imageComponentID, err := context.RestoreUUID(reqTyped.ImageComponentID)
		if err != nil {
			return nil, errors.NewError("画像のコンポーネントIDを作成できません", err)
		}

		sendEmbed, err := send_embed.NewSendEmbed(chID, title, content, colorCode, imageComponentID, reqTyped.DisplayAuthor)
		if err != nil {
			return nil, errors.NewError("Embedアクションを作成できません", err)
		}

		return sendEmbed, nil
	case ReplyEmbedActionReq:
		title, err := action.NewTitle(reqTyped.Title)
		if err != nil {
			return nil, errors.NewError("Titleを作成できません", err)
		}

		content, err := action.NewContent(reqTyped.Content)
		if err != nil {
			return nil, errors.NewError("Contentを作成できません", err)
		}

		colorCode, err := action.NewColorCode(reqTyped.ColorCode)
		if err != nil {
			return nil, errors.NewError("ColorCodeを作成できません", err)
		}

		imageComponentID, err := context.RestoreUUID(reqTyped.ImageComponentID)
		if err != nil {
			return nil, errors.NewError("画像のコンポーネントIDを作成できません", err)
		}

		replyEmbed, err := reply_embed.NewReplyEmbed(
			title,
			content,
			colorCode,
			imageComponentID,
			reqTyped.DisplayAuthor,
			reqTyped.IsEphemeral,
		)
		if err != nil {
			return nil, errors.NewError("Embedアクションを作成できません", err)
		}

		return replyEmbed, nil
	case IfBlockReq:
		condKind, err := condition.NewKind(reqTyped.Condition.Kind)
		if err != nil {
			return nil, errors.NewError("条件の種類を作成できません", err)
		}

		expected, err := condition.NewExpected(reqTyped.Condition.Expected)
		if err != nil {
			return nil, errors.NewError("条件の期待値を作成できません", err)
		}

		cond, err := condition.NewCondition(condKind, expected)
		if err != nil {
			return nil, errors.NewError("条件を作成できません", err)
		}

		// IfBlockReq の場合は、中の`Action`を再帰的に処理する
		var trueActions []action.Action
		for _, reqAct := range reqTyped.TrueAction {
			act, err := CreateActionFromReq(reqAct)
			if err != nil {
				return nil, err
			}
			trueActions = append(trueActions, act)
		}

		var falseActions []action.Action
		for _, reqAct := range reqTyped.FalseAction {
			act, err := CreateActionFromReq(reqAct)
			if err != nil {
				return nil, err
			}
			falseActions = append(falseActions, act)
		}

		ifBlock, err := domain.NewIfBlock(cond, trueActions, falseActions)
		if err != nil {
			return nil, errors.NewError("IfBlockを作成できません", err)
		}

		return ifBlock, nil
	default:
		return nil, errors.NewError("未知のアクションタイプです")
	}
}
