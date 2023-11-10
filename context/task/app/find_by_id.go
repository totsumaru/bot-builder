package app

import (
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// IDでアクションを取得します
func FindByID(tx *gorm.DB, id string) (domain.Task, error) {
	var res domain.Task

	taskID, err := domain.RestoreUUID(id)
	if err != nil {
		return res, errors.NewError("idを復元できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return res, errors.NewError("Gatewayを作成できません", err)
	}

	res, err = gw.FindByID(taskID)
	if err != nil {
		return res, errors.NewError("アクションを取得できません", err)
	}

	return res, nil
}
