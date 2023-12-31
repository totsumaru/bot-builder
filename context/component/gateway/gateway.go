package gateway

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/component/domain"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
	"github.com/totsumaru/bot-builder/context/component/gateway/database"
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

// コンポーネントを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(component domain.Component) error {
	dbComponent, err := castToDBStruct(component)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbComponent)
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
func (g Gateway) Update(component domain.Component) error {
	dbComponent, err := castToDBStruct(component)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&componentDB.Component{}).Where(
		"id = ?",
		dbComponent.ID,
	).Updates(&dbComponent)
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
func (g Gateway) FindByID(id context.UUID) (domain.Component, error) {
	var res domain.Component

	var dbComponent componentDB.Component
	if err := g.tx.First(&dbComponent, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbComponent)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FindByIDs は指定されたIDのリストに基づいてコンポーネントを検索します
func (g Gateway) FindByIDs(ids []context.UUID) ([]domain.Component, error) {
	var dbComponents []componentDB.Component
	var domainComponents []domain.Component

	// IDのリストを文字列のスライスに変換
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = id.String()
	}

	// データベースからIDに一致するレコードを検索
	if err := g.tx.Where("id IN (?)", stringIDs).Find(&dbComponents).Error; err != nil {
		return nil, errors.NewError("複数のIDでコンポーネントを取得できません", err)
	}

	// 取得したDBレコードをドメインモデルに変換
	for _, dbComponent := range dbComponents {
		component, err := castToDomainModel(dbComponent)
		if err != nil {
			return nil, err
		}
		domainComponents = append(domainComponents, component)
	}

	return domainComponents, nil
}

// イベントIDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByEventID(eventID context.UUID) (domain.Component, error) {
	var res domain.Component

	var dbComponent componentDB.Component
	if err := g.tx.First(&dbComponent, "event_id = ?", eventID.String()).Error; err != nil {
		return res, errors.NewError("イベントIDでコンポーネントを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbComponent)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id context.UUID) (domain.Component, error) {
	var res domain.Component

	var dbComponent componentDB.Component
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbComponent, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbComponent)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// 削除します
func (g Gateway) Delete(id context.UUID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&componentDB.Component{}, "id = ?", id.String())
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
func castToDBStruct(component domain.Component) (componentDB.Component, error) {
	dbComponent := componentDB.Component{}

	b, err := json.Marshal(&component)
	if err != nil {
		return dbComponent, errors.NewError("Marshalに失敗しました", err)
	}

	dbComponent.ID = component.ID().String()
	dbComponent.ServerID = component.ServerID().String()
	dbComponent.ApplicationID = component.ApplicationID().String()
	dbComponent.Data = b

	return dbComponent, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbComponent componentDB.Component) (domain.Component, error) {
	m := map[string]any{}
	if err := json.Unmarshal(dbComponent.Data, &m); err != nil {
		return nil, errors.NewError("Unmarshalに失敗しました", err)
	}

	kind := seeker.Str(m, []string{"kind", "value"})
	switch kind {
	case domain.ComponentKindButton:
		res := button.Button{}

		b, err := json.Marshal(m)
		if err != nil {
			return nil, errors.NewError("Marshalに失敗しました", err)
		}

		if err = json.Unmarshal(b, &res); err != nil {
			return nil, errors.NewError("Unmarshalに失敗しました", err)
		}

		return res, nil
	default:
		return nil, errors.NewError("コンポーネントの種類が不正です")
	}
}
