package expose

import (
	"sync"

	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// タスクをキャッシュするためのsync.Mapの利用
var CachedTask = sync.Map{}

// サーバーIDでタスクを取得します
//
// キャッシュが存在しない場合は、DBから取得してキャッシュに追加します。
func GetCachedTasks(serverID string) ([]domain.Task, error) {
	// キャッシュから試みて取得
	if res, ok := CachedTask.Load(serverID); ok {
		return res.([]domain.Task), nil
	}

	// キャッシュにない場合はDBから取得
	tasks, err := app.FindByServerID(DB, serverID)
	if err != nil {
		return nil, errors.NewError("タスクを取得できません", err)
	}

	// 取得したタスクをキャッシュに追加
	CachedTask.Store(serverID, tasks)

	return tasks, nil
}
