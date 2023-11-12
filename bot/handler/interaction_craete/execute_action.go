package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// アクションを実行します
func executeAction(s *discordgo.Session, i *discordgo.InteractionCreate, act action.Action) error {
	switch act.ActionType().String() {
	case action.ActionTypeReplyEmbed:
		replyEmbed := act.(reply_embed.ReplyEmbed)
		embed := &discordgo.MessageEmbed{
			Title:       replyEmbed.Title().String(),
			Description: replyEmbed.Content().String(),
			Color:       replyEmbed.ColorCode().Int(),
		}

		if replyEmbed.DisplayAuthor() {
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    i.Member.User.Username,
				IconURL: i.Member.User.AvatarURL(""),
			}
		}

		componentIDs := make([]string, 0)
		for _, v := range replyEmbed.ComponentID() {
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

		components := make([]discordgo.MessageComponent, 0)
		for _, btn := range btns {
			components = append(components, btn)
		}
		actions := discordgo.ActionsRow{
			Components: components,
		}

		resp := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{actions},
				Embeds:     []*discordgo.MessageEmbed{embed},
				//Flags:      discordgo.MessageFlagsEphemeral, // TODO
			},
		}
		if len(components) == 0 {
			resp.Data.Components = nil
		}

		if err = s.InteractionRespond(i.Interaction, resp); err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	return nil
}
