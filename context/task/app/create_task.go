package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// タスクを作成するリクエストです
type CreateTaskReq struct {
	ServerID      string
	ApplicationID string
	IfBlock       IfBlockReq
}

// タスクを作成します
func CreateTask(tx *gorm.DB, req CreateTaskReq) (domain.Task, error) {
	id, err := context.NewUUID()
	if err != nil {
		return domain.Task{}, errors.NewError("UUIDを作成できません", err)
	}

	serverID, err := context.NewDiscordID(req.ServerID)
	if err != nil {
		return domain.Task{}, errors.NewError("DiscordIDを作成できません", err)
	}

	applicationID, err := context.RestoreUUID(req.ApplicationID)
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

	if err = gw.Create(task); err != nil {
		return domain.Task{}, errors.NewError("タスクを作成できません", err)
	}

	return task, nil
}
