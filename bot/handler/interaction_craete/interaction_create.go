package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	taskApp "github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// インタラクションが作成された時のハンドラです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		domainTasks, err := taskApp.FindByServerID(tx, i.GuildID)
		if err != nil {
			return errors.NewError("タスクを取得できません", err)
		}

		for _, domainTask := range domainTasks {
			kind := domainTask.IfBlock().Condition().Kind().String()
			switch kind {
			case condition.KindClickedButtonIs:
				// ボタンクリックのタイプ以外の場合は無視します
				if i.Type != discordgo.InteractionMessageComponent {
					continue
				}
				// 期待したボタンIDでは無い場合は無視します
				expectedButtonID := domainTask.IfBlock().Condition().Expected().String()
				if i.MessageComponentData().CustomID != expectedButtonID {
					continue
				}

				for _, act := range domainTask.IfBlock().TrueAction() {
					if err = executeAction(s, i, act); err != nil {
						return errors.NewError("処理を実行できません", err)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
		return
	}
}
