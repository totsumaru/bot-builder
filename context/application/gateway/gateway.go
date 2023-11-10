package gateway

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/application/domain"
	"github.com/totsumaru/bot-builder/context/application/gateway/database"
	"github.com/totsumaru/bot-builder/lib/errors"
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

// コンポーネントを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(application domain.Application) error {
	dbApplication, err := castToDBStruct(application)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbApplication)
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
func (g Gateway) Update(application domain.Application) error {
	dbApplication, err := castToDBStruct(application)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&database.Application{}).Where(
		"id = ?",
		dbApplication.ID,
	).Updates(&dbApplication)
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
func (g Gateway) FindByID(id context.UUID) (domain.Application, error) {
	var res domain.Application

	var dbApplication database.Application
	if err := g.tx.First(&dbApplication, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbApplication)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// イベントIDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByEventID(eventID context.UUID) (domain.Application, error) {
	var res domain.Application

	var dbApplication database.Application
	if err := g.tx.First(&dbApplication, "event_id = ?", eventID.String()).Error; err != nil {
		return res, errors.NewError("イベントIDでコンポーネントを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbApplication)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id context.UUID) (domain.Application, error) {
	var res domain.Application

	var dbApplication database.Application
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbApplication, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbApplication)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// 削除します
func (g Gateway) Delete(id context.UUID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&database.Application{}, "id = ?", id.String())
	if result.Error != nil {
		return errors.NewError("削除できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// =============
// private
// =============

// ドメインモデルをDBの構造体に変換します
func castToDBStruct(application domain.Application) (database.Application, error) {
	dbApplication := database.Application{}

	b, err := json.Marshal(&application)
	if err != nil {
		return dbApplication, errors.NewError("Marshalに失敗しました", err)
	}

	dbApplication.ID = application.ID().String()
	dbApplication.ServerID = application.ServerID().String()
	dbApplication.Data = b

	return dbApplication, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbApplication database.Application) (domain.Application, error) {
	var application domain.Application

	if err := json.Unmarshal(dbApplication.Data, &application); err != nil {
		return application, errors.NewError("Unmarshalに失敗しました", err)
	}

	return application, nil
}
