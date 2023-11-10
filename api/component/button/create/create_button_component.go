package create

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/component/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	ServerID      string `json:"server_id"`
	ApplicationID string `json:"application_id"`
	Label         string `json:"label"`
	Style         string `json:"style"`
	URL           string `json:"url"`
}

// レスポンスです
type Res struct {
	ID            string `json:"id"`
	ServerID      string `json:"server_id"`
	ApplicationID string `json:"application_id"`
	Label         string `json:"label"`
	Style         string `json:"style"`
	URL           string `json:"url"`
}

// ボタンコンポーネントを作成します
func CreateButtonComponent(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/component/button/create", func(c *gin.Context) {
		// Tx
		res := Res{}
		err := db.Transaction(func(tx *gorm.DB) error {
			// リクエストをパース
			req := Req{}
			err := c.ShouldBindJSON(&req)
			if err != nil {
				return errors.NewError("リクエストをパースできません", err)
			}

			appReq := app.CreateButtonComponentReq{
				ServerID:      req.ServerID,
				ApplicationID: req.ApplicationID,
				Label:         req.Label,
				Style:         req.Style,
				URL:           req.URL,
			}

			domainButton, err := app.CreateButtonComponent(tx, appReq)
			if err != nil {
				return errors.NewError("ボタンコンポーネントを作成できません", err)
			}

			res.ID = domainButton.ID().String()
			res.ServerID = domainButton.ServerID().String()
			res.ApplicationID = domainButton.ApplicationID().String()
			res.Label = domainButton.Label().String()
			res.Style = domainButton.Style().String()
			res.URL = domainButton.URL().String()

			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"message": "エラーが発生しました"})
			return
		}

		c.JSON(200, res)
	})
}
