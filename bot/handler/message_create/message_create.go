package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/send_text"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		domainTasks, err := app.FindByServerID(tx, m.GuildID)
		if err != nil {
			return errors.NewError("タスクを取得できません", err)
		}

		for _, domainTask := range domainTasks {
			switch domainTask.IfBlock().Condition().Kind().String() {
			case condition.KindCreatedMessageIs:
				if m.Content == domainTask.IfBlock().Condition().Expected().String() {
					for _, act := range domainTask.IfBlock().TrueAction() {
						if err = executeAction(s, m, act); err != nil {
							return errors.NewError("アクションを実行できません", err)
						}
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

// アクションを実行します
func executeAction(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	act action.Action,
) error {
	switch act.ActionType().String() {
	case action.ActionTypeSendText:
		sendText := act.(send_text.SendText)
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			sendText.Content().String(),
		)
		if err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	return nil
}
