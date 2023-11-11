package image

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// 画像コンポーネントです
type Image struct {
	domain.ComponentCore
	url context.URL
}

// 画像コンポーネントを作成します
func NewImage(
	core domain.ComponentCore,
	url context.URL,
) (Image, error) {
	i := Image{
		ComponentCore: core,
		url:           url,
	}

	if err := i.validate(); err != nil {
		return i, err
	}

	return i, nil
}

// 画像のURLを返します
func (i Image) URL() context.URL {
	return i.url
}

// 画像コンポーネントを検証します
func (i Image) validate() error {
	return nil
}

// 画像コンポーネントをJSONに変換します
func (i Image) MarshalJSON() ([]byte, error) {
	data := struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		Kind          domain.Kind       `json:"kind"`
		URL           context.URL       `json:"url"`
	}{
		ID:            i.ID(),
		ServerID:      i.ServerID(),
		ApplicationID: i.ApplicationID(),
		Kind:          i.Kind(),
		URL:           i.url,
	}

	return json.Marshal(data)
}

// JSONから画像コンポーネントを復元します
func (i *Image) UnmarshalJSON(b []byte) error {
	data := struct {
		ID            context.UUID      `json:"id"`
		ServerID      context.DiscordID `json:"server_id"`
		ApplicationID context.UUID      `json:"application_id"`
		Kind          domain.Kind       `json:"kind"`
		URL           context.URL       `json:"url"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	id, err := context.RestoreUUID(data.ID.String())
	if err != nil {
		return errors.NewError("IDを復元できません", err)
	}

	serverID, err := context.NewDiscordID(data.ServerID.String())
	if err != nil {
		return errors.NewError("DiscordIDを作成できません", err)
	}

	appID, err := context.RestoreUUID(data.ApplicationID.String())
	if err != nil {
		return errors.NewError("ApplicationIDを復元できません", err)
	}

	kind, err := domain.NewKind(data.Kind.String())
	if err != nil {
		return errors.NewError("コンポーネントの種類を作成できません", err)
	}

	i.ComponentCore, err = domain.NewComponentCore(id, serverID, appID, kind)
	if err != nil {
		return errors.NewError("コンポーネントの共通の構造体を作成できません", err)
	}

	i.url = data.URL

	return nil
}
