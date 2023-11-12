package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	"github.com/totsumaru/bot-builder/context"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// アクションを実行します
func ExecuteAction(s *discordgo.Session, i *discordgo.InteractionCreate, act action.Action) error {
	switch act.ActionType().String() {
	case action.ActionTypeReplyEmbed:
		replyEmbed, ok := act.(reply_embed.ReplyEmbed)
		if !ok {
			return errors.NewError("型アサーションに失敗しました")
		}
		if err := ReplyEmbedMessage(s, i, replyEmbed); err != nil {
			return errors.NewError("埋め込みメッセージを返信できません", err)
		}
	}

	return nil
}

// 埋め込みメッセージを返信するアクションです
func ReplyEmbedMessage(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	replyEmbed reply_embed.ReplyEmbed,
) error {
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

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			//Flags:      discordgo.MessageFlagsEphemeral, // TODO
		},
	}

	if len(replyEmbed.ComponentID()) > 0 {
		discordBtns, err := GetDiscordButtonsFromComponentIDs(replyEmbed.ComponentID())
		if err != nil {
			return errors.NewError("ボタンを取得できません", err)
		}

		actions := discordgo.ActionsRow{}
		for _, btn := range discordBtns {
			actions.Components = append(actions.Components, btn)
		}
		resp.Data.Components = []discordgo.MessageComponent{actions}
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}

// コンポーネントIDからDiscordのボタンを取得します
func GetDiscordButtonsFromComponentIDs(componentID []context.UUID) ([]discordgo.Button, error) {
	componentIDs := make([]string, 0)
	for _, v := range componentID {
		componentIDs = append(componentIDs, v.String())
	}
	// 複数のIDに一致するコンポーネントを全て取得します
	btnComponents, err := componentApp.FindButtonByIDs(bot.DB, componentIDs)
	if err != nil {
		return nil, errors.NewError("複数IDでコンポーネントを取得できません", err)
	}

	// Discordのボタンに変換します
	btns := make([]discordgo.Button, 0)
	for _, btnComponent := range btnComponents {
		btn := discordgo.Button{
			Label:    btnComponent.Label().String(),
			Style:    bot.ButtonStyleDomainToDiscord[btnComponent.Style().String()],
			CustomID: btnComponent.ID().String(),
		}
		btns = append(btns, btn)
	}

	return btns, nil
}
