package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
	"github.com/totsumaru/bot-builder/context/component/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// 複数のIDでボタンコンポーネントを取得します
func FindButtonByIDs(tx *gorm.DB, id []string) ([]button.Button, error) {
	ids := make([]context.UUID, 0)
	for _, v := range id {
		componentID, err := context.RestoreUUID(v)
		if err != nil {
			return nil, errors.NewError("idを復元できません", err)
		}
		ids = append(ids, componentID)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return nil, errors.NewError("Gatewayを作成できません", err)
	}

	components, err := gw.FindByIDs(ids)
	if err != nil {
		return nil, errors.NewError("コンポーネントを取得できません", err)
	}

	res := make([]button.Button, 0)
	for _, v := range components {
		r, ok := v.(button.Button)
		if !ok {
			return nil, errors.NewError("コンポーネントの型が違います")
		}
		res = append(res, r)
	}

	return res, nil
}
