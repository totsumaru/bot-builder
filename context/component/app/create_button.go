package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
	"github.com/totsumaru/bot-builder/context/component/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// ボタンコンポーネントを作成するリクエストです
type CreateButtonComponentReq struct {
	ServerID      string
	ApplicationID string
	Label         string
	Style         string
	URL           string
}

// ボタンコンポーネントを作成します
func CreateButtonComponent(tx *gorm.DB, req CreateButtonComponentReq) (button.Button, error) {
	id, err := context.NewUUID()
	if err != nil {
		return button.Button{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	serverID, err := context.NewDiscordID(req.ServerID)
	if err != nil {
		return button.Button{}, errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	applicationID, err := context.RestoreUUID(req.ApplicationID)
	if err != nil {
		return button.Button{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	kind, err := domain.NewKind(domain.ComponentKindButton)
	if err != nil {
		return button.Button{}, errors.NewError("コンポーネントの種類の生成に失敗しました", err)
	}

	core, err := domain.NewComponentCore(id, serverID, applicationID, kind)
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

	if err = gw.Create(btn); err != nil {
		return button.Button{}, errors.NewError("ボタンの作成に失敗しました", err)
	}

	return btn, nil
}
