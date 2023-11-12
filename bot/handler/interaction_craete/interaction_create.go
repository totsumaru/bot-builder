package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/bot-builder/bot"
	taskApp "github.com/totsumaru/bot-builder/context/task/app"
	"github.com/totsumaru/bot-builder/context/task/domain"
	"github.com/totsumaru/bot-builder/context/task/domain/condition"
	"github.com/totsumaru/bot-builder/lib/errors"
	"gorm.io/gorm"
)

// インタラクションが作成された時のハンドラです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		domainTasks, err := taskApp.FindByServerID(tx, i.GuildID)
		if err != nil {
			return errors.NewError("タスクを取得できません", err)
		}

		for _, domainTask := range domainTasks {
			kind := domainTask.IfBlock().Condition().Kind().String()
			switch kind {
			case condition.KindClickedButtonIs:
				// ボタンクリックのタイプ以外の場合は無視します
				if i.Type != discordgo.InteractionMessageComponent {
					continue
				}

				if err = CheckAndExecuteActions(s, i, domainTask.IfBlock()); err != nil {
					return errors.NewError("処理を実行できません", err)
				}
			}
		}

		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
		return
	}
}

// 条件を検証し、アクションを実行します
func CheckAndExecuteActions(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	ifBlock domain.IfBlock,
) error {
	// 条件が正しいかどうかを検証します
	ok, err := IsValidCondition(s, i, ifBlock.Condition())
	if err != nil {
		return errors.NewError("条件を検証できません", err)
	}

	actions := ifBlock.FalseAction()
	if ok {
		actions = ifBlock.TrueAction()
	}

	for _, act := range actions {
		if err = ExecuteAction(s, i, act); err != nil {
			return errors.NewError("アクションを実行できません", err)
		}
	}

	return nil
}

// 条件が正しいかどうかを検証します
func IsValidCondition(s *discordgo.Session, i *discordgo.InteractionCreate, cond condition.Condition) (bool, error) {
	expected := cond.Expected().String()

	switch cond.Kind().String() {
	case condition.KindClickedButtonIs:
		if i.MessageComponentData().CustomID == expected {
			return true, nil
		}
		return false, nil
	case condition.KindOperatorIs:
		if i.Member.User.ID == expected {
			return true, nil
		}
		return false, nil
	case condition.KindOperatorRoleHas:
		for _, roleID := range i.Member.Roles {
			if roleID == expected {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.NewError("条件の種類が不正です")
	}
}
