package app

import (
	"github.com/totsumaru/bot-builder/context/action/domain"
	"github.com/totsumaru/bot-builder/context/action/domain/components/button"
	"github.com/totsumaru/bot-builder/context/action/domain/text"
	"github.com/totsumaru/bot-builder/context/action/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// ボタンを作成するリクエストです
type CreateButtonRequest struct {
	Label   string
	Style   string
	EventID string
	URL     string
}

// テキストアクションの作成リクエストです
type CreateTextActionRequest struct {
	EventID     string
	Order       int
	Content     string
	Button      []CreateButtonRequest
	IsResponse  bool
	IsEphemeral bool
	ChannelID   string
}

// テキストアクションを作成します
func NewTextAction(tx *gorm.DB, req CreateTextActionRequest) (text.TextAction, error) {
	res := text.TextAction{}

	eventID, err := domain.RestoreUUID(req.EventID)
	if err != nil {
		return res, errors.NewError("UUIDを復元できません", err)
	}

	order, err := domain.NewOrder(req.Order)
	if err != nil {
		return res, errors.NewError("順序を生成できません", err)
	}

	content, err := text.NewContent(req.Content)
	if err != nil {
		return res, errors.NewError("コンテンツを生成できません", err)
	}

	btns := make([]button.Button, 0)
	for _, b := range req.Button {
		label, err := button.NewLabel(b.Label)
		if err != nil {
			return res, errors.NewError("ラベルを生成できません", err)
		}

		style, err := button.NewStyle(b.Style)
		if err != nil {
			return res, errors.NewError("スタイルを生成できません", err)
		}

		newEventID, err := domain.RestoreUUID(b.EventID)
		if err != nil {
			return res, errors.NewError("UUIDを復元できません", err)
		}

		url, err := button.NewURL(b.URL)
		if err != nil {
			return res, errors.NewError("URLを生成できません", err)
		}

		btn, err := button.NewButton(label, style, newEventID, url)
		if err != nil {
			return res, errors.NewError("ボタンを生成できません", err)
		}

		btns = append(btns, btn)
	}

	channelID, err := domain.NewDiscordID(req.ChannelID)
	if err != nil {
		return res, errors.NewError("DiscordIDを生成できません", err)
	}

	res, err = text.NewTextAction(
		eventID, order, content, btns, req.IsResponse, req.IsEphemeral, channelID,
	)
	if err != nil {
		return res, errors.NewError("テキストアクションを生成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return res, errors.NewError("Gatewayを作成できません", err)
	}

	if err = gw.Create(res); err != nil {
		return res, errors.NewError("テキストアクションを保存できません", err)
	}

	return res, nil
}
