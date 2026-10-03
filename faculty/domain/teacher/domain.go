package teacher

import (
	"context"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/faculty/domain"
	"github.com/violin0622/ink-scroll/faculty/domain/teacher/vo"
)

type Record[ID comparable] interface {
	ID() ID
}

// type TimestampedRecord interface {
// 	GetUpdatedAt() time.Time
// 	GetCreatedAt() time.Time
// 	GetDeletedAt() time.Time
//
// 	SetUpdatedAt(time.Time)
// 	SetCreatedAt(time.Time)
// 	SetDeletedAt(time.Time)
// }

type (
	GetOption    func()
	CreateOption func()
	UpdateOption func()
	DeleteOption func()
)

type Repo[A comparable, B Record[A]] interface {
	Get(context.Context, A, ...GetOption) (B, error)
	Create(context.Context, B) error
	Update(context.Context, B) error
	Delete(context.Context, A) error
}

// type ListOption func()
// type BatchRepo[A comparable, B Record[A]] interface {
// 	List(context.Context, ...ListOption) ([]B, error)
// 	Creates(context.Context, ...B) error
// 	Update(context.Context, ...B) error
// 	Delete(context.Context, ...A) error
// }

// type RelateRepo[A, B comparable, C Record[A], D Record[B]] interface {
// 	ListBy(context.Context, C, ...ListOption) ([]D, error)
// }

//	type DAO interface {
//		Tx(func(tx DAO) error) error
//	}

type TeacherRepo = Repo[ulid.ULID, *Teacher]

// Create implements Repo.
// type TeacherExRepo interface {
// 	Upsert(context.Context, *Teacher) error
// 	GetOrCreate(context.Context, *Teacher) error
// }

type (
	Gender     = vo.Gender
	Birthday   = vo.Birthday
	HireStatus = vo.HireStatus
)

type (
	EventRepo = domain.EventRepo[ulid.ULID, ulid.ULID]
	Aggregate = domain.Aggregate[ulid.ULID]
)
