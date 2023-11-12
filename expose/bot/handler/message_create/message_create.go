package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/expose"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := expose.DB.Transaction(func(tx *gorm.DB) error {
		domainTasks, err := expose.GetCachedTasks(m.GuildID)
		if err != nil {
			return errors.NewError("タスクを取得できません", err)
		}

		for _, domainTask := range domainTasks {
			if err = CheckAndExecuteActions(s, m, domainTask.IfBlock()); err != nil {
				return errors.NewError("アクションを実行できません", err)
			}
		}

		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
		return
	}
}

// 条件を検証し、アクションを実行します
func CheckAndExecuteActions(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	ifBlock domain.IfBlock,
) error {
	// 条件が正しいかどうかを検証します
	ok, err := IsValidCondition(s, m, ifBlock.Condition())
	if err != nil {
		return errors.NewError("条件を検証できません", err)
	}

	actions := ifBlock.FalseAction()
	if ok {
		actions = ifBlock.TrueAction()
	}

	for _, act := range actions {
		if err = ExecuteAction(s, m, act); err != nil {
			return errors.NewError("アクションを実行できません", err)
		}
	}

	return nil
}

// 条件が正しいかどうかを検証します
func IsValidCondition(s *discordgo.Session, m *discordgo.MessageCreate, cond condition.Condition) (bool, error) {
	expected := cond.Expected().String()

	switch cond.Kind().String() {
	case condition.KindCreatedMessageIs:
		if m.Content == expected {
			return true, nil
		}
		return false, nil
	case condition.KindOperatorIs:
		if m.Author.ID == expected {
			return true, nil
		}
		return false, nil
	case condition.KindOperatorRoleHas:
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			return false, errors.NewError("メンバーを取得できません", err)
		}
		for _, roleID := range member.Roles {
			if roleID == expected {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.NewError("条件の種類が不正です")
	}
}
