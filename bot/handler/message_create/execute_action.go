package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
	"github.com/totsumaru/bot-builder/context/task/domain/action/send_text"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// アクションを実行します
func executeAction(s *discordgo.Session, m *discordgo.MessageCreate, act action.Action) error {
	switch act.ActionType().String() {
	case action.ActionTypeSendText:
		sendText := act.(send_text.SendText)

		componentIDs := make([]string, 0)
		for _, v := range sendText.ComponentID() {
			componentIDs = append(componentIDs, v.String())
		}
		btnComponents, err := componentApp.FindButtonByIDs(bot.DB, componentIDs)
		if err != nil {
			return errors.NewError("複数IDでコンポーネントを取得できません", err)
		}

		btns := make([]discordgo.Button, 0)
		for _, btnComponent := range btnComponents {
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

		_, err = s.ChannelMessageSendComplex(sendText.ChannelID().String(), data)
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
