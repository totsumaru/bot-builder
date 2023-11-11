package create

// IfBlockのリクエストです
type IfBlockReq struct {
	ActionType string `json:"action_type"`
	Condition  struct {
		Kind     string `json:"kind"`
		Expected string `json:"expected"`
	} `json:"condition"`
	TrueAction  []map[string]any `json:"true_action"`
	FalseAction []map[string]any `json:"false_action"`
}

// ===============================================================
// Actionのリクエスト
//
// コードは使用されませんが、リクエストのために定義しておきます
// ===============================================================

// テキストを送信するアクションのリクエストです
type SendTextActionReq struct {
	ActionType  string   `json:"action_type"`
	ChannelID   string   `json:"channel_id"`
	Content     string   `json:"content"`
	ComponentID []string `json:"component_id"`
}

// テキストを返信するアクションのリクエストです
type ReplyTextActionReq struct {
	ActionType  string   `json:"action_type"`
	Content     string   `json:"content"`
	IsEphemeral bool     `json:"is_ephemeral"`
	ComponentID []string `json:"component_id"`
}

// Embedを送信するアクションのリクエストです
type SendEmbedActionReq struct {
	ActionType       string `json:"action_type"`
	ChannelID        string `json:"channel_id"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	ColorCode        int    `json:"color_code"`
	ImageComponentID string `json:"image_component_id"`
	DisplayAuthor    bool   `json:"display_author"`
}

// Embedを返信するアクションのリクエストです
type ReplyEmbedActionReq struct {
	ActionType       string `json:"action_type"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	ColorCode        int    `json:"color_code"`
	ImageComponentID string `json:"image_component_id"`
	DisplayAuthor    bool   `json:"display_author"`
	IsEphemeral      bool   `json:"is_ephemeral"`
}
