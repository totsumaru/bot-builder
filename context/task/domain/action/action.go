package action

// アクションのインターフェースです
type Action interface {
	ActionType() ActionType
}
