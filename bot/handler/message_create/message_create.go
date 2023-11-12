package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	taskApp "github.com/totsumaru/bot-builder/context/task/app"
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

// TODO: 優先度低: FindByIDsでボタンを複数一気に取れるようにする
// アクションを実行します
func executeAction(s *discordgo.Session, m *discordgo.MessageCreate, act action.Action) error {
	switch act.ActionType().String() {
	case action.ActionTypeSendText:
		sendText := act.(send_text.SendText)

		btns := make([]discordgo.Button, 0)
		for _, componentID := range sendText.ComponentID() {
			btnComponent, err := componentApp.FindButtonByID(bot.DB, componentID.String())
			if err != nil {
				return errors.NewError("コンポーネントを取得できません", err)
			}
			btn := discordgo.Button{
				Label:    btnComponent.Label().String(),
				Style:    bot.ButtonStyleDomainToDiscord[btnComponent.Style().String()],
				CustomID: btnComponent.ID().String(),
			}
			btns = append(btns, btn)
		}

		actions := discordgo.ActionsRow{}
		for _, btn := range btns {
			actions.Components = append(actions.Components, btn)
		}

		data := &discordgo.MessageSend{
			Content:    sendText.Content().String(),
			Components: []discordgo.MessageComponent{actions},
		}

		_, err := s.ChannelMessageSendComplex(sendText.ChannelID().String(), data)
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

		btns := make([]discordgo.Button, 0)
		for _, componentID := range replyEmbed.ComponentID() {
			btnComponent, err := componentApp.FindButtonByID(bot.DB, componentID.String())
			if err != nil {
				return errors.NewError("コンポーネントを取得できません", err)
			}
			btn := discordgo.Button{
				Label:    btnComponent.Label().String(),
				Style:    bot.ButtonStyleDomainToDiscord[btnComponent.Style().String()],
				CustomID: btnComponent.ID().String(),
			}
			btns = append(btns, btn)
		}

		actions := discordgo.ActionsRow{}
		for _, btn := range btns {
			actions.Components = append(actions.Components, btn)
		}

		data := &discordgo.MessageSend{
			Embed:      embed,
			Components: []discordgo.MessageComponent{actions},
			Reference:  m.Reference(),
		}

		_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
		if err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	default:
		return errors.NewError("アクションタイプが存在していません")
	}

	return nil
}
