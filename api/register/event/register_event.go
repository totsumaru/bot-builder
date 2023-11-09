package event

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/event/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// イベントを作成します
func RegisterEvent(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/register/event", func(c *gin.Context) {
		// Tx
		err := db.Transaction(func(tx *gorm.DB) error {
			req := app.CreateMessageEventRequest{
				AllowedRoleID:    c.PostFormArray("allowed_role_id"),
				AllowedChannelID: c.PostFormArray("allowed_channel_id"),
				Keyword:          c.PostForm("keyword"),
				MatchType:        c.PostForm("match_type"),
			}
			_, err := app.CreateMessageEvent(tx, req)
			if err != nil {
				return errors.NewError("イベントを作成できません", err)
			}

			return nil
		})
		if err != nil {
			fmt.Println(errors.NewError("エラーが発生しました", err))
			return
		}

		c.JSON(200, nil)
	})
}
