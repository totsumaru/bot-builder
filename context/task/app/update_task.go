package app

import (
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// タスクを更新します
func UpdateTextTask(tx *gorm.DB, req UpdateTaskReq) (domain.Task, error) {
	id, err := domain.RestoreUUID(req.ID)
	if err != nil {
		return domain.Task{}, errors.NewError("UUIDを復元できません", err)
	}

	serverID, err := domain.NewDiscordID(req.ServerID)
	if err != nil {
		return domain.Task{}, errors.NewError("DiscordIDを作成できません", err)
	}

	applicationID, err := domain.RestoreUUID(req.ApplicationID)
	if err != nil {
		return domain.Task{}, errors.NewError("UUIDを作成できません", err)
	}

	ifBlock, err := CreateIfBlock(req.IfBlock)
	if err != nil {
		return domain.Task{}, errors.NewError("IfBlockを作成できません", err)
	}

	task, err := domain.NewTask(id, serverID, applicationID, ifBlock)
	if err != nil {
		return domain.Task{}, errors.NewError("タスクを作成できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return domain.Task{}, errors.NewError("Gatewayを作成できません", err)
	}

	if err = gw.Update(task); err != nil {
		return domain.Task{}, errors.NewError("タスクを更新できません", err)
	}

	return task, nil
}
