package name

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/application/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// アプリケーション名を更新します
func UpdateApplicationName(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/application/update/name", func(c *gin.Context) {
		// Tx
		err := db.Transaction(func(tx *gorm.DB) error {
			// リクエストをパース
			req := Req{}
			if err := c.ShouldBindJSON(&req); err != nil {
				return errors.NewError("リクエストをパースできません", err)
			}

			if err := app.UpdateName(tx, req.ID, req.Name); err != nil {
				return errors.NewError("アプリケーション名を更新できません", err)
			}

			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"message": "エラーが発生しました"})
			return
		}

		c.JSON(200, nil)
	})
}
