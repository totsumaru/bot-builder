package create

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/lib/errors"
	"github.com/totsumaru/bot-builder/lib/seeker"
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
			if err := c.ShouldBindJSON(&req); err != nil {
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
			fmt.Println(err)
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
func castApiActionReqToAppActionReq(apiAction map[string]any) (any, error) {
	actionType := seeker.Str(apiAction, []string{"action_type"})
	switch actionType {
	case ActionTypeSendText:
		appReq := app.SendTextActionReq{
			ChannelID:   seeker.Str(apiAction, []string{"channel_id"}),
			Content:     seeker.Str(apiAction, []string{"content"}),
			ComponentID: seeker.SliceString(apiAction, []string{"component_id"}),
		}

		return appReq, nil
	case ActionTypeReplyText:
		appReq := app.ReplyTextActionReq{
			Content:     seeker.Str(apiAction, []string{"content"}),
			IsEphemeral: seeker.Bool(apiAction, []string{"is_ephemeral"}),
			ComponentID: seeker.SliceString(apiAction, []string{"component_id"}),
		}
		return appReq, nil
	case ActionTypeSendEmbed:
		appReq := app.SendEmbedActionReq{
			ChannelID:        seeker.Str(apiAction, []string{"channel_id"}),
			Title:            seeker.Str(apiAction, []string{"title"}),
			Content:          seeker.Str(apiAction, []string{"content"}),
			ColorCode:        seeker.Int(apiAction, []string{"color_code"}),
			ImageComponentID: seeker.Str(apiAction, []string{"image_component_id"}),
			DisplayAuthor:    seeker.Bool(apiAction, []string{"display_author"}),
		}
		return appReq, nil
	case ActionTypeReplyEmbed:
		appReq := app.ReplyEmbedActionReq{
			Title:            seeker.Str(apiAction, []string{"title"}),
			Content:          seeker.Str(apiAction, []string{"content"}),
			ColorCode:        seeker.Int(apiAction, []string{"color_code"}),
			ImageComponentID: seeker.Str(apiAction, []string{"image_component_id"}),
			DisplayAuthor:    seeker.Bool(apiAction, []string{"display_author"}),
			IsEphemeral:      seeker.Bool(apiAction, []string{"is_ephemeral"}),
		}
		return appReq, nil
	case ActionTypeIfBlock:
		trueAction := make([]any, 0)
		for _, apiReqTrueAction := range seeker.Slice(apiAction, []string{"true_action"}) {
			appReq, err := castApiActionReqToAppActionReq(apiReqTrueAction)
			if err != nil {
				return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
			}
			trueAction = append(trueAction, appReq)
		}

		falseAction := make([]any, 0)
		for _, apiReqFalseAction := range seeker.Slice(apiAction, []string{"false_action"}) {
			appReq, err := castApiActionReqToAppActionReq(apiReqFalseAction)
			if err != nil {
				return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
			}
			falseAction = append(falseAction, appReq)
		}

		appReq := app.IfBlockReq{}
		appReq.Condition.Kind = seeker.Str(apiAction, []string{"condition", "kind"})
		appReq.Condition.Expected = seeker.Str(apiAction, []string{"condition", "expected"})
		appReq.TrueAction = trueAction
		appReq.FalseAction = falseAction

		return appReq, nil
	default:
		return app.IfBlockReq{}, errors.NewError("Actionの型を変換できません", nil)
	}
}
