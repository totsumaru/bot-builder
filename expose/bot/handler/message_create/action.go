package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/context"
	componentApp "github.com/totsumaru/bot-builder/context/component/app"
	taskDomain "github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/action"
	"github.com/totsumaru/bot-builder/context/task/domain/action/reply_embed"
	"github.com/totsumaru/bot-builder/context/task/domain/action/send_text"
	"github.com/totsumaru/bot-builder/expose"
	"github.com/totsumaru/bot-builder/expose/bot"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// アクションを実行します
func ExecuteAction(s *discordgo.Session, m *discordgo.MessageCreate, act action.Action) error {
	switch act.ActionType().String() {
	case action.ActionTypeSendText:
		sendText, ok := act.(send_text.SendText)
		if !ok {
			return errors.NewError("型アサーションに失敗しました")
		}
		if err := SendTextMessage(s, sendText); err != nil {
			return errors.NewError("テキストメッセージを送信できません", err)
		}
	case action.ActionTypeReplyEmbed:
		replyEmbed, ok := act.(reply_embed.ReplyEmbed)
		if !ok {
			return errors.NewError("型アサーションに失敗しました")
		}
		if err := ReplyEmbedMessage(s, m, replyEmbed); err != nil {
			return errors.NewError("埋め込みメッセージを返信できません", err)
		}
	case action.ActionTypeIfBlock:
		ifBlock, ok := act.(taskDomain.IfBlock)
		if !ok {
			return errors.NewError("型アサーションに失敗しました")
		}
		// ifBlockの場合は、再帰的にアクションを実行します
		if err := CheckAndExecuteActions(s, m, ifBlock); err != nil {
			return errors.NewError("アクションを実行できません", err)
		}
	default:
		return errors.NewError("アクションタイプが存在していません")
	}

	return nil
}

// メッセージを送信するアクションです
func SendTextMessage(s *discordgo.Session, sendText send_text.SendText) error {
	data := &discordgo.MessageSend{
		Content: sendText.Content().String(),
	}

	if len(sendText.ComponentID()) > 0 {
		discordBtns, err := GetDiscordButtonsFromComponentIDs(sendText.ComponentID())
		if err != nil {
			return errors.NewError("ボタンを取得できません", err)
		}

		actions := discordgo.ActionsRow{}
		for _, btn := range discordBtns {
			actions.Components = append(actions.Components, btn)
		}
		data.Components = []discordgo.MessageComponent{actions}
	}

	_, err := s.ChannelMessageSendComplex(sendText.ChannelID().String(), data)
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}

// 埋め込みメッセージを返信するアクションです
func ReplyEmbedMessage(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	replyEmbed reply_embed.ReplyEmbed,
) error {
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

	data := &discordgo.MessageSend{
		Embed:     embed,
		Reference: m.Reference(),
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
		data.Components = []discordgo.MessageComponent{actions}
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
	if err != nil {
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
	btnComponents, err := componentApp.FindButtonByIDs(expose.DB, componentIDs)
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
