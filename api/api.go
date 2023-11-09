package api

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/api/register/action"
	"github.com/totsumaru/bot-builder/api/register/event"
	"gorm.io/gorm"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine, db *gorm.DB) {
	Route(e)
	event.RegisterEvent(e, db)
	action.RegisterAction(e, db)
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
