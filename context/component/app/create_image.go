package app

import (
	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/context/component/domain/image"
	"github.com/totsumaru/bot-builder/context/component/gateway"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// 画像コンポーネントを作成するリクエストです
type CreateImageComponentReq struct {
	ServerID      string
	ApplicationID string
	URL           string
}

// 画像コンポーネントを作成します
func CreateImageComponent(tx *gorm.DB, req CreateImageComponentReq) (image.Image, error) {
	id, err := context.NewUUID()
	if err != nil {
		return image.Image{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	serverID, err := context.NewDiscordID(req.ServerID)
	if err != nil {
		return image.Image{}, errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	applicationID, err := context.RestoreUUID(req.ApplicationID)
	if err != nil {
		return image.Image{}, errors.NewError("UUIDの生成に失敗しました", err)
	}

	kind, err := domain.NewKind(domain.ComponentKindImage)
	if err != nil {
		return image.Image{}, errors.NewError("コンポーネントの種類の生成に失敗しました", err)
	}

	core, err := domain.NewComponentCore(id, serverID, applicationID, kind)
	if err != nil {
		return image.Image{}, errors.NewError("コンポーネントの共通部分の生成に失敗しました", err)
	}

	url, err := context.NewURL(req.URL)
	if err != nil {
		return image.Image{}, errors.NewError("URLの生成に失敗しました", err)
	}

	img, err := image.NewImage(core, url)
	if err != nil {
		return image.Image{}, errors.NewError("画像の生成に失敗しました", err)
	}

	gw, err := gateway.NewGateway(tx)
	if err != nil {
		return image.Image{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	if err = gw.Create(img); err != nil {
		return image.Image{}, errors.NewError("画像の作成に失敗しました", err)
	}

	return img, nil
}
