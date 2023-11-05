package app

import (
	"github.com/totsumaru/bot-builder/context/event/domain"
	"github.com/totsumaru/bot-builder/context/event/domain/message"
	"github.com/totsumaru/bot-builder/context/event/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// メッセージ作成のイベントを作成するリクエストです
type CreateMessageEventRequest struct {
	allowedRoleID    []string
	allowedChannelID []string
	keyword          string
	matchType        string
}

// メッセージ作成のイベントを作成します
func NewMessageEvent(tx *gorm.DB, req CreateMessageEventRequest) (message.MessageEvent, error) {
	res := message.MessageEvent{}

	allowedRoleID := make([]domain.DiscordID, 0)
	for _, id := range req.allowedRoleID {
		dID, err := domain.NewDiscordID(id)
		if err != nil {
			return res, errors.NewError("DiscordIDを生成できません", err)
		}
		allowedRoleID = append(allowedRoleID, dID)
	}

	allowedChannelID := make([]domain.DiscordID, 0)
	for _, id := range req.allowedChannelID {
		dID, err := domain.NewDiscordID(id)
		if err != nil {
			return res, errors.NewError("DiscordIDを生成できません", err)
		}
		allowedChannelID = append(allowedChannelID, dID)
	}

	keyword, err := message.NewKeyword(req.keyword)
	if err != nil {
		return res, errors.NewError("キーワードを生成できません", err)
	}

	matchType, err := message.NewMatchType(req.matchType)
	if err != nil {
		return res, errors.NewError("一致条件を生成できません", err)
	}

	res, err = message.NewMessageEvent(allowedRoleID, allowedChannelID, keyword, matchType)
	if err != nil {
		return res, errors.NewError("メッセージイベントを生成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return res, errors.NewError("Gatewayを作成できません", err)
	}

	if err = gw.Create(res); err != nil {
		return res, errors.NewError("メッセージイベントを保存できません", err)
	}

	return res, nil
}
