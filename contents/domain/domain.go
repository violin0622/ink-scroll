package domain

import (
	"context"

	"github.com/violin0622/ink-scroll/common"
)

type EventRepo[EID, AID comparable] interface {
	// Fetch(context.Context, AID) ([]Event[EID, AID], error)
	Append(context.Context, ...common.Event[EID, AID]) error
}
