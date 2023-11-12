package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	taskApp "github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		domainTasks, err := taskApp.FindByServerID(tx, m.GuildID)
		if err != nil {
			return errors.NewError("タスクを取得できません", err)
		}

		for _, domainTask := range domainTasks {
			// `送られたメッセージがxだったら` の条件でなければ、実行されません
			kind := domainTask.IfBlock().Condition().Kind().String()
			if kind != condition.KindCreatedMessageIs {
				continue
			}

			expectedText := domainTask.IfBlock().Condition().Expected().String()
			if m.Content == expectedText {
				for _, act := range domainTask.IfBlock().TrueAction() {
					if err = executeAction(s, m, act); err != nil {
						return errors.NewError("trueアクションを実行できません", err)
					}
				}
			} else {
				for _, act := range domainTask.IfBlock().FalseAction() {
					if err = executeAction(s, m, act); err != nil {
						return errors.NewError("falseアクションを実行できません", err)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
		return
	}
}
