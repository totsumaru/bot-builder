package action

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/action/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	EventID     string `json:"event_id"`
	Order       int    `json:"order"`
	Content     string `json:"content"`
	IsResponse  bool   `json:"is_response"`
	IsEphemeral bool   `json:"is_ephemeral"`
	ChannelID   string `json:"channel_id"`
	Button      []struct {
		Label string `json:"label"`
		Style string `json:"style"`
		URL   string `json:"url"`
	} `json:"button"`
}

// アクションを作成します
func RegisterAction(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/register/action", func(c *gin.Context) {
		// Tx
		err := db.Transaction(func(tx *gorm.DB) error {
			// リクエストをパース
			req := Req{}
			err := c.ShouldBindJSON(&req)
			if err != nil {
				return errors.NewError("リクエストをパースできません", err)
			}

			// ボタンのイベントを保存します
			{

			}

			// ボタンのリクエストをパース
			btns := make([]app.CreateButtonRequest, 0)
			for _, b := range req.Button {
				btn := app.CreateButtonRequest{
					Label: b.Label,
					Style: b.Style,
					URL:   b.URL,
				}
				btns = append(btns, btn)
			}

			actionAppReq := app.CreateTextActionRequest{
				EventID:     req.EventID,
				Order:       req.Order,
				Content:     req.Content,
				IsResponse:  req.IsResponse,
				IsEphemeral: req.IsEphemeral,
				ChannelID:   req.ChannelID,
				Button:      btns,
			}
			_, err = app.CreateTextAction(tx, actionAppReq)
			if err != nil {
				return errors.NewError("テキストアクションを作成できません", err)
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
