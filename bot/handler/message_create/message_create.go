package message_create

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	"github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
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
	fmt.Println("アクションタイプ: ", act.ActionType().String())

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
	case action.ActionTypeReplyEmbed:
		replyEmbed := act.(reply_embed.ReplyEmbed)
		embed := &discordgo.MessageEmbed{
			Title:       replyEmbed.Title().String(),
			Description: replyEmbed.Content().String(),
			Color:       replyEmbed.ColorCode().Int(),
		}

		if replyEmbed.DisplayAuthor() {
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			}
		}

		_, err := s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
		if err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	default:

	}

	return nil
}
