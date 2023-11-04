package trigger

import (
	"github.com/bwmarrin/discordgo"
)

// ボタントリガーです
type ButtonTrigger struct {
	id string
}

// トリガーのIDを返します
func (t ButtonTrigger) ID() string {
	return t.id
}

// トリガーがマッチするかどうかを返します
func (t ButtonTrigger) IsMatch(i interface{}) (bool, error) {
	interactionCreate, ok := i.(*discordgo.InteractionCreate)
	if !ok {
		return false, nil
	}

	if interactionCreate.Type != discordgo.InteractionMessageComponent {
		return false, nil
	}

	if interactionCreate.MessageComponentData().CustomID != t.id {
		return false, nil
	}

	return true, nil
}
