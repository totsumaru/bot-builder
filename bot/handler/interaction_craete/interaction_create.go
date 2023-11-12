package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	taskApp "github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
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
		panic(err)
	}
}

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

		if err := s.InteractionRespond(i.Interaction, resp); err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	return nil
}
