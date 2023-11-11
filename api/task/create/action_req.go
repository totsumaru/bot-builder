package create

const (
	ActionTypeIfBlock    = "if_block"
	ActionTypeSendText   = "send_text"
	ActionTypeReplyText  = "reply_text"
	ActionTypeSendEmbed  = "send_embed"
	ActionTypeReplyEmbed = "reply_embed"
)

type Action interface {
	ActionTypeString() string
}

// IfBlockのリクエストです
type IfBlockReq struct {
	ActionType string `json:"action_type"`
	Condition  struct {
		Kind     string `json:"kind"`
		Expected string `json:"expected"`
	} `json:"condition"`
	TrueAction  []Action `json:"true_action"`
	FalseAction []Action `json:"false_action"`
}

func (req IfBlockReq) ActionTypeString() string {
	return req.ActionType
}

// ==============================================
// Actionのリクエスト
// ==============================================

// テキストを送信するアクションのリクエストです
type SendTextActionReq struct {
	ActionType  string   `json:"action_type"`
	ChannelID   string   `json:"channel_id"`
	Content     string   `json:"content"`
	ComponentID []string `json:"component_id"`
}

func (req SendTextActionReq) ActionTypeString() string {
	return req.ActionType
}

// テキストを返信するアクションのリクエストです
type ReplyTextActionReq struct {
	ActionType  string   `json:"action_type"`
	Content     string   `json:"content"`
	IsEphemeral bool     `json:"is_ephemeral"`
	ComponentID []string `json:"component_id"`
}

func (req ReplyTextActionReq) ActionTypeString() string {
	return req.ActionType
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

func (req SendEmbedActionReq) ActionTypeString() string {
	return req.ActionType
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

func (req ReplyEmbedActionReq) ActionTypeString() string {
	return req.ActionType
}
