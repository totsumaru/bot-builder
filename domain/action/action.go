package action

import "github.com/totsumaru/bot-builder/domain"

// アクションのInterfaceです
type Action interface {
	ID() domain.UUID
}
