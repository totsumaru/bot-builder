package create

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// リクエストです
type Req struct {
	ServerID      string     `json:"server_id"`
	ApplicationID string     `json:"application_id"`
	IfBlock       IfBlockReq `json:"if_block"`
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

			appIfBlock, err := castIfBlockReqToAppReq(req.IfBlock)
			if err != nil {
				return errors.NewError("IfBlockの型を変換できません", err)
			}

			appReq := app.CreateTaskReq{
				ServerID:      req.ServerID,
				ApplicationID: req.ApplicationID,
				IfBlock:       appIfBlock,
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

// IfBlockのリクエストをAppのリクエストに変換します
//
// ActionをAppの型に変換します。
func castIfBlockReqToAppReq(apiIfBlockReq IfBlockReq) (app.IfBlockReq, error) {
	res := app.IfBlockReq{}
	res.Condition.Kind = apiIfBlockReq.Condition.Kind
	res.Condition.Expected = apiIfBlockReq.Condition.Expected

	// APIリクエストをAppの型に変換します
	trueActions := make([]any, 0)
	for _, apiReqTrueAction := range apiIfBlockReq.TrueAction {
		appReq, err := castApiActionReqToAppActionReq(apiReqTrueAction)
		if err != nil {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		trueActions = append(trueActions, appReq)
	}

	falseActions := make([]any, 0)
	for _, apiReqFalseAction := range apiIfBlockReq.FalseAction {
		appReq, err := castApiActionReqToAppActionReq(apiReqFalseAction)
		if err != nil {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		falseActions = append(falseActions, appReq)
	}

	res.TrueAction = trueActions
	res.FalseAction = falseActions

	return res, nil
}

// ActionのリクエストをAppのリクエストに変換します
func castApiActionReqToAppActionReq(apiAction Action) (any, error) {
	// any型をActionの型に変換します
	switch apiAction.ActionTypeString() {
	case ActionTypeSendText:
		apiActionReq, ok := apiAction.(SendTextActionReq)
		if !ok {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		appReq := app.SendTextActionReq{
			ChannelID:   apiActionReq.ChannelID,
			Content:     apiActionReq.Content,
			ComponentID: apiActionReq.ComponentID,
		}
		return appReq, nil
	case ActionTypeReplyText:
		apiActionReq, ok := apiAction.(ReplyTextActionReq)
		if !ok {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		appReq := app.ReplyTextActionReq{
			Content:     apiActionReq.Content,
			IsEphemeral: apiActionReq.IsEphemeral,
			ComponentID: apiActionReq.ComponentID,
		}
		return appReq, nil
	case ActionTypeSendEmbed:
		apiActionReq, ok := apiAction.(SendEmbedActionReq)
		if !ok {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		appReq := app.SendEmbedActionReq{
			ChannelID:        apiActionReq.ChannelID,
			Title:            apiActionReq.Title,
			Content:          apiActionReq.Content,
			ColorCode:        apiActionReq.ColorCode,
			ImageComponentID: apiActionReq.ImageComponentID,
			DisplayAuthor:    apiActionReq.DisplayAuthor,
		}
		return appReq, nil
	case ActionTypeReplyEmbed:
		apiActionReq, ok := apiAction.(ReplyEmbedActionReq)
		if !ok {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}
		appReq := app.ReplyEmbedActionReq{
			Title:            apiActionReq.Title,
			Content:          apiActionReq.Content,
			ColorCode:        apiActionReq.ColorCode,
			ImageComponentID: apiActionReq.ImageComponentID,
			DisplayAuthor:    apiActionReq.DisplayAuthor,
			IsEphemeral:      apiActionReq.IsEphemeral,
		}
		return appReq, nil
	case ActionTypeIfBlock:
		apiActionReq, ok := apiAction.(IfBlockReq)
		if !ok {
			return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
		}

		trueAction := make([]any, 0)
		for _, apiReqTrueAction := range apiActionReq.TrueAction {
			appReq, err := castApiActionReqToAppActionReq(apiReqTrueAction)
			if err != nil {
				return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
			}
			trueAction = append(trueAction, appReq)
		}

		falseAction := make([]any, 0)
		for _, apiReqFalseAction := range apiActionReq.FalseAction {
			appReq, err := castApiActionReqToAppActionReq(apiReqFalseAction)
			if err != nil {
				return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
			}
			falseAction = append(falseAction, appReq)
		}

		appReq := app.IfBlockReq{}
		appReq.Condition.Kind = apiActionReq.Condition.Kind
		appReq.Condition.Expected = apiActionReq.Condition.Expected
		appReq.TrueAction = trueAction
		appReq.FalseAction = falseAction

		return appReq, nil
	default:
		return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
	}
}
