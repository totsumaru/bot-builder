package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// サーバーIDでタスクを取得します
func FindByServerID(tx *gorm.DB, serverID string) ([]domain.Task, error) {
	res := make([]domain.Task, 0)

	sID, err := context.NewDiscordID(serverID)
	if err != nil {
		return res, errors.NewError("サーバーIDを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return res, errors.NewError("Gatewayを作成できません", err)
	}

	res, err = gw.FindByServerID(sID)
	if err != nil {
		return res, errors.NewError("サーバーIDでタスクを取得できません", err)
	}

	return res, nil
}
