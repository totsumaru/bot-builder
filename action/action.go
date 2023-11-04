package action

import "github.com/bwmarrin/discordgo"

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
			TextSendAction{
				id:        "action1",
				triggerID: "trigger1",
				order:     1,
				text:      "Hello!",
				channelID: "", // テスト
			},
		},
	}

	return ActionsDB[triggerID]
}
