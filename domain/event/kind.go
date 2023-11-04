package event

// Eventの種類を表す定数です
type Kind string

const (
	// メッセージが送信された時のイベントです
	EventKindMessageCreate Kind = "MESSAGE"
	// ボタンが押された時のイベントです
	EventKindButton Kind = "BUTTON"
)
