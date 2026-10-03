package novel

import (
	"context"
	"errors"

	"github.com/puzpuzpuz/xsync/v4"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

var _ novel.Repo = (*repo)(nil)

type repo struct {
	m *xsync.Map[novel.ID, novel.Novel]
}

func New() repo {
	return repo{m: xsync.NewMap[novel.ID, novel.Novel]()}
}

func Newx(m *xsync.Map[novel.ID, novel.Novel]) repo {
	return repo{m: m}
}

// Get implements domain.AggregateRepo.
func (r repo) Get(_ context.Context, id novel.ID) (novel.Novel, error) {
	if n, ok := r.m.Load(id); ok {
		return n, nil
	}

	return novel.Novel{}, novel.NewNotfoundError(errors.New(`novel not found`))
}

// Save implements domain.AggregateRepo.
func (r repo) Save(_ context.Context, n novel.Novel) error {
	r.m.Store(n.ID(), n)
	return nil
}

func (r repo) Update(
	ctx context.Context,
	id novel.ID,
	fn func(context.Context, *novel.Novel) error,
) (err error) {
	r.m.Compute(
		id,
		func(n novel.Novel, loaded bool) (
			novel.Novel, xsync.ComputeOp,
		) {
			if !loaded {
				err = novel.NewNotfoundError(errors.New(`novel not found`))
				return novel.Novel{}, xsync.CancelOp
			}

			if err = fn(ctx, &n); err != nil {
				return novel.Novel{}, xsync.CancelOp
			}

			return n, xsync.UpdateOp
		})

	return err
}

// var _ novel.Repo = (*esrepo)(nil)

// type esrepo struct {
// 	f novel.EventRepo
// }
//
// func NewEventSourceRepo(es novel.EventRepo) esrepo {
// 	return esrepo{f: fake.New[novel.ID, novel.ID]()}
// }
//
// func (e esrepo) Get(ctx context.Context, id novel.ID) (novel.Novel, error) {
// 	es, err := e.f.Fetch(ctx, id)
// 	if err != nil {
// 		return novel.Novel{}, fmt.Errorf(`get: %w`, err)
// 	}
//
// 	n := novel.New(e.f)
// 	if err = n.ReplayEvents(es...); err != nil {
// 		return novel.Novel{}, fmt.Errorf(`get: %w`, err)
// 	}
//
// 	return n, nil
// }
//
// func (esrepo) Save(context.Context, novel.Novel) error { return nil }
//
// func (e esrepo) Update(
// 	context.Context,
// 	novel.ID,
// 	func(context.Context, *novel.Novel) error,
// ) error {
// 	return nil
// }
