package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/application/domain"
	"github.com/totsumaru/bot-builder/context/application/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// 名前を変更します
func UpdateName(tx *gorm.DB, id, name string) error {
	i, err := context.RestoreUUID(id)
	if err != nil {
		return errors.NewError("UUIDの復元に失敗しました", err)
	}

	n, err := domain.NewName(name)
	if err != nil {
		return errors.NewError("名前の生成に失敗しました", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	application, err := gw.FindByIDForUpdate(i)
	if err != nil {
		return errors.NewError("アプリケーションの取得に失敗しました", err)
	}

	if err = application.UpdateName(n); err != nil {
		return errors.NewError("名前の変更に失敗しました", err)
	}

	if err = gw.Update(application); err != nil {
		return errors.NewError("アプリケーションの更新に失敗しました", err)
	}

	return nil
}
