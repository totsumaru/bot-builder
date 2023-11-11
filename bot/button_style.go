package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/context/component/domain/button"
)

// ドメインのボタンスタイルと、Discordのボタンスタイルのmappingです
var ButtonStyleDomainToDiscord = map[string]discordgo.ButtonStyle{
	button.ButtonStylePrimary:   discordgo.PrimaryButton,
	button.ButtonStyleSecondary: discordgo.SecondaryButton,
	button.ButtonStyleSuccess:   discordgo.SuccessButton,
	button.ButtonStyleDanger:    discordgo.DangerButton,
	button.ButtonStyleLink:      discordgo.LinkButton,
}
