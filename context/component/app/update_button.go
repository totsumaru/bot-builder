package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
	"github.com/totsumaru/bot-builder/context/component/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// ボタンコンポーネントを更新するリクエストです
type UpdateButtonComponentReq struct {
	ID            string
	ServerID      string
	ApplicationID string
	Label         string
	Style         string
	URL           string
}

// ボタンコンポーネントを更新します
func UpdateButtonComponent(tx *gorm.DB, req UpdateButtonComponentReq) (button.Button, error) {
	id, err := context.RestoreUUID(req.ID)
	if err != nil {
		return button.Button{}, errors.NewError("UUIDを復元できません", err)
	}

	serverID, err := context.NewDiscordID(req.ServerID)
	if err != nil {
		return button.Button{}, errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	applicationID, err := context.RestoreUUID(req.ApplicationID)
	if err != nil {
		return button.Button{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	core, err := domain.NewComponentCore(id, serverID, applicationID)
	if err != nil {
		return button.Button{}, errors.NewError("コンポーネントの共通部分の生成に失敗しました", err)
	}

	label, err := button.NewLabel(req.Label)
	if err != nil {
		return button.Button{}, errors.NewError("ラベルの生成に失敗しました", err)
	}

	style, err := button.NewStyle(req.Style)
	if err != nil {
		return button.Button{}, errors.NewError("スタイルの生成に失敗しました", err)
	}

	url, err := context.NewURL(req.URL)
	if err != nil {
		return button.Button{}, errors.NewError("URLの生成に失敗しました", err)
	}

	btn, err := button.NewButton(core, label, style, url)
	if err != nil {
		return button.Button{}, errors.NewError("ボタンの生成に失敗しました", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return button.Button{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	if err = gw.Update(btn); err != nil {
		return button.Button{}, errors.NewError("ボタンの更新に失敗しました", err)
	}

	return btn, nil
}
