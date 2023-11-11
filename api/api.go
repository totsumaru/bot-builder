package api

import (
	"github.com/gin-gonic/gin"
	applicationCreate "github.com/totsumaru/bot-builder/api/application/create"
	"github.com/totsumaru/bot-builder/api/application/update/name"
	buttonCreate "github.com/totsumaru/bot-builder/api/component/button/create"
	"github.com/totsumaru/bot-builder/api/component/button/update"
	"github.com/totsumaru/bot-builder/api/task/create"
	"gorm.io/gorm"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine, db *gorm.DB) {
	Route(e)
	// application
	applicationCreate.CreateApplication(e, db)
	name.UpdateApplicationName(e, db)
	// component
	buttonCreate.CreateButtonComponent(e, db)
	update.UpdateButtonComponent(e, db)
	// task
	create.CreateTask(e, db)
}

// ルートです
//
// Note: この関数は削除しても問題ありません
func Route(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
