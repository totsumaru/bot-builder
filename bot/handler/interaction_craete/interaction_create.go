package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot/action"
	"github.com/totsumaru/bot-builder/bot/trigger"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// インタラクションが作成された時のハンドラです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, t := range trigger.Triggers[i.GuildID] {
		match, err := t.IsMatch(i)
		if err != nil {
			errors.SendErrMsg(s, err)
			continue
		}
		if match {
			// アクションを実行する(複数実行する場合もあるのでループ)
			for _, act := range action.GetActions(t.ID()) {
				if err = act.Execute(s, i); err != nil {
					errors.SendErrMsg(s, err)
					continue
				}
			}
		}
	}
}
