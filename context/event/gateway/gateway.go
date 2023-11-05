package gateway

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/event/domain"
	"github.com/totsumaru/bot-builder/context/event/domain/message"
	"github.com/totsumaru/bot-builder/context/event/gateway/database"
	"github.com/totsumaru/bot-builder/lib/errors"
	"github.com/totsumaru/bot-builder/lib/seeker"
	"gorm.io/gorm"
)

type Gateway struct {
	tx *gorm.DB
}

// gatewayを作成します
func NewGateway(tx *gorm.DB) (Gateway, error) {
	if tx == nil {
		return Gateway{}, errors.NewError("引数が空です")
	}

	res := Gateway{
		tx: tx,
	}

	return res, nil
}

// イベントを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(e domain.Event) error {
	dbEvent, err := castToDBStruct(e)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbEvent)
	if result.Error != nil {
		return errors.NewError("レコードを保存できませんでした", result.Error)
	}
	// 主キー制約違反を検出（同じIDのレコードが既に存在する場合）
	if result.RowsAffected == 0 {
		return errors.NewError("既存のレコードが存在しています")
	}

	return nil
}

// 更新します
func (g Gateway) Update(u domain.Event) error {
	dbEvent, err := castToDBStruct(u)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&database.Event{}).Where(
		"id = ?",
		dbEvent.ID,
	).Updates(&dbEvent)
	if result.Error != nil {
		return errors.NewError("更新できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// IDでイベントを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id domain.UUID) (domain.Event, error) {
	var res domain.Event

	var dbEvent database.Event
	if err := g.tx.First(&dbEvent, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでイベントを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbEvent)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでイベントを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id domain.UUID) (domain.Event, error) {
	var res domain.Event

	var dbEvent database.Event
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbEvent, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでイベントを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbEvent)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// =============
// private
// =============

// ドメインモデルをDBの構造体に変換します
func castToDBStruct(eventInterface domain.Event) (database.Event, error) {
	res := database.Event{}

	b, err := json.Marshal(&eventInterface)
	if err != nil {
		return res, errors.NewError("Marshalに失敗しました", err)
	}

	res.ID = eventInterface.ID().String()
	res.Data = b

	return res, nil
}

// DBの構造体からドメインモデルに変換します
//
// kindの種類によって、返す構造体が変わります。
func castToDomainModel(dbEvent database.Event) (domain.Event, error) {
	// 一度、map[string]interface{}に変換します
	m := map[string]interface{}{}
	if err := json.Unmarshal(dbEvent.Data, &m); err != nil {
		return nil, errors.NewError("Unmarshalに失敗しました", err)
	}

	// kindの種類によって、返す構造体が変わります
	kind := seeker.Str(m, []string{"kind", "value"})
	switch kind {
	case domain.EventKindMessageCreate:
		res := message.MessageEvent{}
		if err := json.Unmarshal(dbEvent.Data, &res); err != nil {
			return nil, errors.NewError("MessageEventでUnmarshalに失敗しました", err)
		}
		return res, nil
	}

	return nil, errors.NewError("不明なイベントの種類です")
}
