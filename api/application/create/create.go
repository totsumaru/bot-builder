package create

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/application/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	ServerID string `json:"server_id"`
	Name     string `json:"name"`
}

// レスポンスです
type Res struct {
	ID       string `json:"id"`
	ServerID string `json:"server_id"`
	Name     string `json:"name"`
}

// アプリケーションを作成します
func CreateApplication(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/application/create", func(c *gin.Context) {
		// Tx
		res := Res{}
		err := db.Transaction(func(tx *gorm.DB) error {
			// リクエストをパース
			req := Req{}
			err := c.ShouldBindJSON(&req)
			if err != nil {
				return errors.NewError("リクエストをパースできません", err)
			}

			domainApplication, err := app.CreateApplication(tx, req.ServerID, req.Name)
			if err != nil {
				return errors.NewError("アプリケーションを作成できません", err)
			}

			res.ID = domainApplication.ID().String()
			res.ServerID = domainApplication.ServerID().String()
			res.Name = domainApplication.Name().String()

			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"message": "エラーが発生しました"})
			return
		}

		c.JSON(200, res)
	})
}
