package trigger

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/lib/channel"
	"github.com/totsumaru/bot-builder/lib/errors"
	"github.com/totsumaru/bot-builder/lib/role"
)

const (
	MatchTypeComplete = "complete" // 完全一致
	MatchTypePartial  = "partial"  // 部分一致
)

// メッセージ作成のトリガーです
type MessageTrigger struct {
	id        string
	keyword   string
	matchType string // 完全一致 or 部分一致
	allow     struct {
		roleID    []string
		channelID []string
	}
}

// トリガーのIDを返します
func (t MessageTrigger) ID() string {
	return t.id
}

// トリガーがマッチするかどうかを返します
func (t MessageTrigger) IsMatch(m interface{}) (bool, error) {
	messageCreate, ok := m.(*discordgo.MessageCreate)
	if !ok {
		return false, nil
	}

	// ロールIDが設定されている場合はロールIDが一致するか確認する
	// 許可されたロールを持っていない場合は、falseを返す
	if len(t.allow.roleID) > 0 {
		if !role.HasAllowedRoleID(t.allow.roleID, messageCreate.Member.Roles) {
			return false, nil
		}
	}

	// チャンネルIDが設定されている場合はチャンネルIDが一致するか確認する
	// 許可されたチャンネルではない場合は、falseを返す
	if len(t.allow.channelID) > 0 {
		if !channel.IsAllowedChannelID(t.allow.channelID, messageCreate.ChannelID) {
			return false, nil
		}
	}

	switch t.matchType {
	case MatchTypeComplete:
		return t.keyword == messageCreate.Content, nil
	case MatchTypePartial:
		return strings.Contains(messageCreate.Content, t.keyword), nil
	}

	return false, errors.NewError("MessageTriggerのmatchTypeが不正です")
}
