package gateway

import (
	"encoding/json"

	"github.com/totsumaru/bot-builder/context"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/gateway/database"
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

// タスクを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(task domain.Task) error {
	dbTask, err := castToDBStruct(task)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbTask)
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
func (g Gateway) Update(u domain.Task) error {
	dbTask, err := castToDBStruct(u)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&database.Task{}).Where(
		"id = ?",
		dbTask.ID,
	).Updates(&dbTask)
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
func (g Gateway) FindByID(id context.UUID) (domain.Task, error) {
	var res domain.Task

	var dbTask database.Task
	if err := g.tx.First(&dbTask, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbTask)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// イベントIDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByServerID(serverID context.DiscordID) ([]domain.Task, error) {
	res := make([]domain.Task, 0)

	dbTask := make([]database.Task, 0)
	if err := g.tx.Find(&dbTask, "server_id = ?", serverID.String()).Error; err != nil {
		return res, errors.NewError("サーバーIDでTaskを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	for _, v := range dbTask {
		r, err := castToDomainModel(v)
		if err != nil {
			return res, errors.NewError("DBをドメインモデルに変換できません", err)
		}

		res = append(res, r)
	}

	return res, nil
}

// FOR UPDATEでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id context.UUID) (domain.Task, error) {
	var res domain.Task

	var dbTask database.Task
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbTask, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbTask)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// 削除します
func (g Gateway) Delete(id context.UUID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&database.Task{}, "id = ?", id.String())
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
func castToDBStruct(task domain.Task) (database.Task, error) {
	dbTask := database.Task{}

	b, err := json.Marshal(&task)
	if err != nil {
		return dbTask, errors.NewError("Marshalに失敗しました", err)
	}

	dbTask.ID = task.ID().String()
	dbTask.ServerID = task.ServerID().String()
	dbTask.ApplicationID = task.ApplicationID().String()
	dbTask.Data = b

	return dbTask, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbTask database.Task) (domain.Task, error) {
	task := domain.Task{}

	if err := json.Unmarshal(dbTask.Data, &task); err != nil {
		return task, errors.NewError("Unmarshalに失敗しました", err)
	}

	return task, nil
}
