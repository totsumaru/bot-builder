package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/application/domain"
	"github.com/totsumaru/bot-builder/context/application/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// アプリケーションを作成します
func CreateApplication(tx *gorm.DB, serverID string, name string) (domain.Application, error) {
	id, err := context.NewUUID()
	if err != nil {
		return domain.Application{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	sid, err := context.NewDiscordID(serverID)
	if err != nil {
		return domain.Application{}, errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	n, err := domain.NewName(name)
	if err != nil {
		return domain.Application{}, errors.NewError("アプリケーション名の生成に失敗しました", err)
	}

	application, err := domain.NewApplication(id, sid, n)
	if err != nil {
		return domain.Application{}, errors.NewError("アプリケーションの生成に失敗しました", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return domain.Application{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	if err = gw.Create(application); err != nil {
		return domain.Application{}, errors.NewError("アプリケーションの作成に失敗しました", err)
	}

	return application, nil
}
