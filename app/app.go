package app

import (
	"encoding/json"
	"fmt"

	"github.com/totsumaru/bot-builder/domain"
	"github.com/totsumaru/bot-builder/domain/event/message"
	"github.com/totsumaru/bot-builder/lib/errors"
)

// メッセージイベントを作成します
func CreateMessageEvent(
	allowedRoleID []string,
	allowedChannelID []string,
	keyword string,
	matchType string,
) error {
	allowedRoles := make([]domain.DiscordID, 0)
	for _, id := range allowedRoleID {
		allowedRole, err := domain.NewDiscordID(id)
		if err != nil {
			return errors.NewError("DiscordIDを生成できません", err)
		}
		allowedRoles = append(allowedRoles, allowedRole)
	}

	allowedChannels := make([]domain.DiscordID, 0)
	for _, id := range allowedChannelID {
		allowedChannel, err := domain.NewDiscordID(id)
		if err != nil {
			return errors.NewError("DiscordIDを生成できません", err)
		}
		allowedChannels = append(allowedChannels, allowedChannel)
	}

	kw, err := message.NewKeyword(keyword)
	if err != nil {
		return errors.NewError("キーワードを生成できません", err)
	}

	mt, err := message.NewMatchType(matchType)
	if err != nil {
		return errors.NewError("一致条件を生成できません", err)
	}

	me, err := message.NewMessageEvent(allowedRoles, allowedChannels, kw, mt)
	if err != nil {
		return errors.NewError("メッセージイベントを生成できません", err)
	}

	b, err := json.Marshal(me)
	if err != nil {
		return errors.NewError("JSONに変換できません", err)
	}

	me2 := message.MessageEvent{}
	if err = json.Unmarshal(b, &me2); err != nil {
		return errors.NewError("JSONから構造体に変換できません", err)
	}

	fmt.Printf("%+v\n", me2)

	// DBに保存します

	return nil
}
