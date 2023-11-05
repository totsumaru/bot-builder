package gateway

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context/action/domain"
	"github.com/totsumaru/bot-builder/context/action/domain/text"
	"github.com/totsumaru/bot-builder/context/action/gateway/database"
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

// アクションを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(act domain.Action) error {
	dbAction, err := castToDBStruct(act)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbAction)
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
func (g Gateway) Update(u domain.Action) error {
	dbAction, err := castToDBStruct(u)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&database.Action{}).Where(
		"id = ?",
		dbAction.ID,
	).Updates(&dbAction)
	if result.Error != nil {
		return errors.NewError("更新できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// IDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id domain.UUID) (domain.Action, error) {
	var res domain.Action

	var dbAction database.Action
	if err := g.tx.First(&dbAction, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbAction)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// イベントIDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByEventID(eventID domain.UUID) (domain.Action, error) {
	var res domain.Action

	var dbAction database.Action
	if err := g.tx.First(&dbAction, "event_id = ?", eventID.String()).Error; err != nil {
		return res, errors.NewError("イベントIDでActionを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbAction)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id domain.UUID) (domain.Action, error) {
	var res domain.Action

	var dbAction database.Action
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbAction, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbAction)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// =============
// private
// =============

// ドメインモデルをDBの構造体に変換します
func castToDBStruct(actionInterface domain.Action) (database.Action, error) {
	res := database.Action{}

	b, err := json.Marshal(&actionInterface)
	if err != nil {
		return res, errors.NewError("Marshalに失敗しました", err)
	}

	res.ID = actionInterface.ID().String()
	res.EventID = actionInterface.EventID().String()
	res.Data = b

	return res, nil
}

// DBの構造体からドメインモデルに変換します
//
// kindの種類によって、返す構造体が変わります。
func castToDomainModel(dbAction database.Action) (domain.Action, error) {
	// 一度、map[string]interface{}に変換します
	m := map[string]interface{}{}
	if err := json.Unmarshal(dbAction.Data, &m); err != nil {
		return nil, errors.NewError("Unmarshalに失敗しました", err)
	}

	// kindの種類によって、返す構造体が変わります
	kind := seeker.Str(m, []string{"kind", "value"})
	switch kind {
	case domain.ActionKindText:
		res := text.TextAction{}
		if err := json.Unmarshal(dbAction.Data, &res); err != nil {
			return nil, errors.NewError("TextActionでUnmarshalに失敗しました", err)
		}
		return res, nil
	}

	return nil, errors.NewError("不明なActionの種類です")
}
