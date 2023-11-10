package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// タスクを削除します
func DeleteTask(tx *gorm.DB, id string) error {
	i, err := context.RestoreUUID(id)
	if err != nil {
		return errors.NewError("idを復元できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return errors.NewError("Gatewayを作成できません", err)
	}

	if err = gw.Delete(i); err != nil {
		return errors.NewError("タスクを削除できません", err)
	}

	return nil
}
