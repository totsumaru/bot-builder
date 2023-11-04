package action

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// テキストを送信するアクションです
type TextSendAction struct {
	id        string
	triggerID string
	order     int
	// ここから独自のフィールド
	text      string
	channelID string // 空の場合は、メッセージを送信したチャンネルに送信する
}

// アクションのIDを返します
func (a TextSendAction) ID() string {
	return a.id
}

// トリガーIDを返します
func (a TextSendAction) TriggerID() string {
	return a.triggerID
}

// アクションの実行順を返します
func (a TextSendAction) Order() int {
	return a.order
}

// アクションを実行します
func (a TextSendAction) Execute(s *discordgo.Session, event interface{}) error {
	channelID := a.channelID

	// チャンネルが指定されていない場合は、送信元のチャンネルに送信する
	if channelID == "" {
		switch event.(type) {
		case *discordgo.MessageCreate:
			channelID = event.(*discordgo.MessageCreate).ChannelID
		case *discordgo.InteractionCreate:
			channelID = event.(*discordgo.InteractionCreate).ChannelID
		}
	}

	if _, err := s.ChannelMessageSend(channelID, a.text); err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
