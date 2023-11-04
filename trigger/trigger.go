package trigger

// トリガーのInterfaceです
type Trigger interface {
	IsMatch(event interface{}) (bool, error)
	ID() string
}

// サーバー名: トリガーの配列
var Triggers = map[string][]Trigger{
	"984614055681613864": {
		MessageTrigger{
			id:        "trigger1",
			keyword:   "hello",
			matchType: "complete",
			allow: struct {
				roleID    []string
				channelID []string
			}{
				roleID:    []string{"998800967665459240"},
				channelID: []string{}},
		},
	},
}
