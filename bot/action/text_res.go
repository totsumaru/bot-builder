package action

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot/action/components/button"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// 返信のテキストを送信するアクションです
type TextResAction struct {
	id        string
	triggerID string
	order     int
	// ここから独自のフィールド
	text        string
	button      []button.Button
	isEphemeral bool
}

// アクションのIDを返します
func (a TextResAction) ID() string {
	return a.id
}

// トリガーIDを返します
func (a TextResAction) TriggerID() string {
	return a.triggerID
}

// アクションの実行順を返します
func (a TextResAction) Order() int {
	return a.order
}

// アクションを実行します
func (a TextResAction) Execute(s *discordgo.Session, event interface{}) error {
	switch event.(type) {
	case *discordgo.MessageCreate:
		m := event.(*discordgo.MessageCreate)
		btns := make([]discordgo.MessageComponent, 0)

		for _, b := range a.button {
			btn := discordgo.Button{
				Label:    b.Label,
				Style:    a.ButtonStyle(b),
				CustomID: b.TriggerID,
			}

			btns = append(btns, btn)
		}

		// ActionRowの作成
		actionRow := discordgo.ActionsRow{
			Components: btns,
		}

		messageSend := &discordgo.MessageSend{
			Content: a.text,
			Components: []discordgo.MessageComponent{
				actionRow,
			},
			Reference: m.Reference(),
		}

		if _, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend); err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	case *discordgo.InteractionCreate:
		i := event.(*discordgo.InteractionCreate)

		resp := &discordgo.InteractionResponse{
			//Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: a.text,
			},
		}

		if a.isEphemeral {
			resp.Data.Flags = discordgo.MessageFlagsEphemeral
		}

		if err := s.InteractionRespond(i.Interaction, resp); err != nil {
			return errors.NewError("メッセージを送信できません", err)
		}
	}

	return nil
}

// ボタンのスタイルを返します
func (a TextResAction) ButtonStyle(b button.Button) discordgo.ButtonStyle {
	var style discordgo.ButtonStyle

	switch b.Style {
	case button.ButtonStylePrimary:
		style = discordgo.PrimaryButton
	case button.ButtonStyleSecondary:
		style = discordgo.SecondaryButton
	case button.ButtonStyleSuccess:
		style = discordgo.SuccessButton
	case button.ButtonStyleDanger:
		style = discordgo.DangerButton
	case button.ButtonStyleLink:
		style = discordgo.LinkButton
	}

	return style
}
