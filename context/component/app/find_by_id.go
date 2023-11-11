package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
	"github.com/totsumaru/bot-builder/context/component/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// IDでボタンコンポーネントを取得します
func FindButtonByID(tx *gorm.DB, id string) (button.Button, error) {
	componentID, err := context.RestoreUUID(id)
	if err != nil {
		return button.Button{}, errors.NewError("idを復元できません", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return button.Button{}, errors.NewError("Gatewayを作成できません", err)
	}

	componentRes, err := gw.FindByID(componentID)
	if err != nil {
		return button.Button{}, errors.NewError("コンポーネントを取得できません", err)
	}

	buttonComponent, ok := componentRes.(button.Button)
	if !ok {
		return button.Button{}, errors.NewError("コンポーネントの型が違います")
	}

	return buttonComponent, nil
}
