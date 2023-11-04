package action

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot/action/components/button"
)

// アクションのInterfaceです
type Action interface {
	Execute(s *discordgo.Session, event interface{}) error
	ID() string
	TriggerID() string
	Order() int
}

// トリガーIDからアクションを取得します
func GetActions(triggerID string) []Action {
	// トリガーIDに対応するアクションを返す
	var ActionsDB = map[string][]Action{
		"trigger1": {
			TextResAction{
				id:        "action1",
				triggerID: "trigger1",
				order:     2,
				text:      "これは返信です",
				button: []button.Button{
					{
						Label:     "ボタン1",
						Style:     button.ButtonStylePrimary,
						URL:       "",
						Emoji:     "",
						TriggerID: "trigger2",
					},
				},
				isEphemeral: true,
			},
		},
		"trigger2": {
			TextResAction{
				id:          "action2",
				triggerID:   "trigger2",
				order:       1,
				text:        "Hello!",
				isEphemeral: true,
			},
		},
	}

	return ActionsDB[triggerID]
}
