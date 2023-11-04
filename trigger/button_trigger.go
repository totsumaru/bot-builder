package trigger

// ボタントリガーです
type ButtonTrigger struct {
	id string
}

// トリガーのIDを返します
func (t ButtonTrigger) ID() string {
	return t.id
}

// トリガーがマッチするかどうかを返します
func (t ButtonTrigger) IsMatch(m interface{}) (bool, error) {
	return false, nil
}
