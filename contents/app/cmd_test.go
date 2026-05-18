package app

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"play-ddd/contents/app/command"
	"play-ddd/contents/domain/novel"
	"play-ddd/utils/retry"
)

var _ = Describe("CommandHandler", func() {
	Describe("CreateNovel", func() {
		var (
			ctx     context.Context
			ctrl    *gomock.Controller
			repo    *MockRepo
			tx      *MockRepo
			novels  *MockNovelRepo
			handler CommandHandler
			cmd     command.CreateNovel
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl = gomock.NewController(GinkgoT())
			repo = NewMockRepo(ctrl)
			tx = NewMockRepo(ctrl)
			novels = NewMockNovelRepo(ctrl)
			handler = NewCommandHandler(repo, logr.Discard())
			cmd = command.CreateNovel{
				AuthorID: ulid.Make(),
				Title:    "An awesome novel",
				Desc:     "从前有座山山上有座庙庙里有个老和尚...",
				Category: "志异",
				Tags:     []string{"志怪"},
			}
		})

		It("does not retry ordinary transaction errors", func() {
			ordinaryErr := errors.New("ordinary failure")
			repo.EXPECT().
				Tx(gomock.Any()).
				Return(ordinaryErr).
				Times(1)

			_, err := handler.CreateNovel(ctx, cmd)

			Expect(err).To(MatchError(ContainSubstring("create novel")))
			Expect(errors.Is(err, ordinaryErr)).To(BeTrue())
		})

		It("retries transient transaction errors", func() {
			transientErr := errors.New("temporary failure")
			attempts := 0
			repo.EXPECT().
				Tx(gomock.Any()).
				DoAndReturn(func(func(Repo) error) error {
					attempts++
					if attempts < 3 {
						return retry.Transient(transientErr)
					}
					return nil
				}).
				Times(3)

			_, err := handler.CreateNovel(ctx, cmd)

			Expect(err).NotTo(HaveOccurred())
			Expect(attempts).To(Equal(3))
		})

		It("rejects a title already used by the same author", func() {
			repo.EXPECT().
				Tx(gomock.Any()).
				DoAndReturn(func(fn func(Repo) error) error {
					return fn(tx)
				})
			tx.EXPECT().
				Novel().
				Return(novels)
			novels.EXPECT().
				ByAuthorAndTitle(gomock.Any(), cmd.AuthorID, cmd.Title).
				Return(novel.Novel{}, nil)

			_, err := handler.CreateNovel(ctx, cmd)

			Expect(errors.Is(err, novel.ErrTitleAlreadyExist)).To(BeTrue())
		})

		It("creates and saves a novel when the title is available", func() {
			events := &recordingEventRepo{}
			var saved novel.Novel

			repo.EXPECT().
				Tx(gomock.Any()).
				DoAndReturn(func(fn func(Repo) error) error {
					return fn(tx)
				})
			tx.EXPECT().
				Novel().
				Return(novels)
			novels.EXPECT().
				ByAuthorAndTitle(gomock.Any(), cmd.AuthorID, cmd.Title).
				Return(novel.Novel{}, novel.NewNotfoundError(errors.New("not found")))
			tx.EXPECT().
				Event().
				Return(events)
			tx.EXPECT().
				Novel().
				Return(novels)
			novels.EXPECT().
				Save(gomock.Any(), gomock.AssignableToTypeOf(novel.Novel{})).
				DoAndReturn(func(_ context.Context, n novel.Novel) error {
					saved = n
					return nil
				})

			id, err := handler.CreateNovel(ctx, cmd)

			Expect(err).NotTo(HaveOccurred())
			Expect(id).NotTo(Equal(novel.ZeroID))
			Expect(saved.ID()).To(Equal(id))
			Expect(events.appendCalls).To(Equal(1))
			Expect(events.appended).To(HaveLen(1))
		})
	})
})

type recordingEventRepo struct {
	appendCalls int
	appended    []novel.Event
}

func (r *recordingEventRepo) Append(_ context.Context, evs ...novel.Event) error {
	r.appendCalls++
	r.appended = append(r.appended, evs...)
	return nil
}
