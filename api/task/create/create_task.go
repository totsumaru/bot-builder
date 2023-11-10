package create

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	ServerID      string `json:"server_id"`
	ApplicationID string `json:"application_id"`
	IfBlock       struct {
		Condition struct {
			Kind     string `json:"kind"`
			Expected string `json:"expected"`
		} `json:"condition"`
		TrueAction  []any `json:"true_action"`
		FalseAction []any `json:"false_action"`
	}
}

// レスポンスです
type Res struct {
	ID string `json:"id"`
}

// タスクを作成します
func CreateTask(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/task/create", func(c *gin.Context) {
		// Tx
		res := Res{}
		err := db.Transaction(func(tx *gorm.DB) error {
			// リクエストをパース
			req := Req{}
			err := c.ShouldBindJSON(&req)
			if err != nil {
				return errors.NewError("リクエストをパースできません", err)
			}

			appReq := app.CreateTaskReq{
				ServerID:      req.ServerID,
				ApplicationID: req.ApplicationID,
				IfBlock:       app.IfBlockReq(req.IfBlock),
			}

			domainTask, err := app.CreateTask(tx, appReq)
			if err != nil {
				return errors.NewError("タスクを作成できません", err)
			}

			res.ID = domainTask.ID().String()

			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"message": "エラーが発生しました"})
			return
		}

		c.JSON(200, res)
	})
}
