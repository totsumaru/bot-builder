package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/action"
	"github.com/totsumaru/bot-builder/lib/errors"
	"github.com/totsumaru/bot-builder/trigger"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// トリガーでマッチしたらアクションを実行する
	for _, t := range trigger.Triggers[m.GuildID] {
		match, err := t.IsMatch(m)
		if err != nil {
			errors.SendErrMsg(s, err)
			continue
		}
		if match {
			// アクションを実行する(複数実行する場合もあるのでループ)
			for _, act := range action.GetActions(t.ID()) {
				if err = act.Execute(s, m); err != nil {
					errors.SendErrMsg(s, err)
					continue
				}
			}
		}
	}
}
