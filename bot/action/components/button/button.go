package button

const (
	ButtonStylePrimary   = "PRIMARY"
	ButtonStyleSecondary = "SECONDARY"
	ButtonStyleSuccess   = "SUCCESS"
	ButtonStyleDanger    = "DANGER"
	ButtonStyleLink      = "LINK"
)

// ボタンのコンポーネントです
type Button struct {
	Label     string
	Style     string
	URL       string
	Emoji     string
	TriggerID string
}
