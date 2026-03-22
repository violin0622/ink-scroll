package outbox_test

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"play-ddd/contents/infra/outbox"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

var _ outbox.EventBus = &mockEventBus{}

type mockEventBus struct {
	published int
}

// Pub implements [outbox.EventBus].
func (eb *mockEventBus) Pub(_ context.Context, evs []*novelv1.Event) error {
	eb.published += len(evs)
	return nil
}

// func ids(events []*novelv1.Event) string {
// 	var s strings.Builder
// 	for i := range events {
// 		s.WriteString(events[i].GetId().String() + ` `)
// 	}
// 	return s.String()
// }

var _ outbox.EventRepo = &mockEventRepo{}

type mockEventRepo struct {
	ev     []*novelv1.Event
	cursor int
}

func newRepo() *mockEventRepo {
	return &mockEventRepo{ev: events}
}

// Process implements [outbox.EventRepo].
func (r *mockEventRepo) Process(
	ctx context.Context,
	fn func(outbox.EventsBatch) error,
) error {
	eb := &mockEventbatch{r: r, cursor: r.cursor}
	if err := fn(eb); err != nil {
		return err
	}
	r.cursor = eb.cursor
	return nil
}

var _ outbox.EventsBatch = &mockEventbatch{}

type mockEventbatch struct {
	r      *mockEventRepo
	cursor int // absolute index
}

// Poll implements [outbox.EventsBatch].
func (b *mockEventbatch) Poll(a outbox.Arg) ([]*novelv1.Event, error) {
	right := min(len(b.r.ev), b.cursor+a.Limit)
	return b.r.ev[b.cursor:right], nil
}

// AdvanceCursor implements [outbox.EventsBatch].
func (b *mockEventbatch) AdvanceCursor(rs ...outbox.Result) error {
	b.cursor += len(rs)
	return nil
}

// test data
var (
	aggKindA, aggKindB                 = `kA`, `kB`
	eventKindA, eventKindB, eventKindC = `eA`, `eB`, `eC`

	eventsIDs = []ulid.ULID{
		ulid.Make(),
		ulid.Make(),
		ulid.Make(),
		ulid.Make(),
		ulid.Make(),
	}
	aggregateIDs = []ulid.ULID{
		ulid.Make(),
		ulid.Make(),
		ulid.Make(),
	}

	events = []*novelv1.Event{{
		Id:            ulidpb.From(eventsIDs[0]),
		AggregateId:   ulidpb.From(aggregateIDs[0]),
		EmitAt:        timestamppb.New(time.Now()),
		Kind:          eventKindA,
		AggregateKind: aggKindA,
	}, {
		Id:            ulidpb.From(eventsIDs[1]),
		AggregateId:   ulidpb.From(aggregateIDs[1]),
		EmitAt:        timestamppb.New(time.Now()),
		Kind:          eventKindA,
		AggregateKind: aggKindA,
	}, {
		Id:            ulidpb.From(eventsIDs[2]),
		AggregateId:   ulidpb.From(aggregateIDs[0]),
		EmitAt:        timestamppb.New(time.Now()),
		Kind:          eventKindB,
		AggregateKind: aggKindA,
	}, {
		Id:            ulidpb.From(eventsIDs[3]),
		AggregateId:   ulidpb.From(aggregateIDs[2]),
		EmitAt:        timestamppb.New(time.Now()),
		Kind:          eventKindC,
		AggregateKind: aggKindB,
	}, {
		Id:            ulidpb.From(eventsIDs[4]),
		AggregateId:   ulidpb.From(aggregateIDs[1]),
		EmitAt:        timestamppb.New(time.Now()),
		Kind:          eventKindB,
		AggregateKind: aggKindA,
	}}
)

func Example_mockEventRepoProcessMultipleTimes() {
	repo := newRepo()
	err := repo.Process(context.TODO(), func(eb outbox.EventsBatch) error {
		evs, err := eb.Poll(outbox.Arg{Limit: 2})
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(len(evs))
		fmt.Println(evs[0].GetId().Into() == events[0].GetId().Into())
		fmt.Println(evs[1].GetId().Into() == events[1].GetId().Into())

		return eb.AdvanceCursor(
			outbox.Result{
				ID:     evs[0].GetId().Into(),
				Status: outbox.StatusPublished,
			},
			outbox.Result{
				ID:     evs[1].GetId().Into(),
				Status: outbox.StatusPublished,
			},
		)
	})
	fmt.Println(err)

	err = repo.Process(context.TODO(), func(eb outbox.EventsBatch) error {
		evs, err := eb.Poll(outbox.Arg{Limit: 2})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// events 0-1 was processed, this time should got 2-3
		fmt.Println(len(evs))
		fmt.Println(evs[0].GetId().Into() == events[2].GetId().Into())
		fmt.Println(evs[1].GetId().Into() == events[3].GetId().Into())

		return eb.AdvanceCursor(
			outbox.Result{
				ID:     evs[0].GetId().Into(),
				Status: outbox.StatusPublished,
			},
			outbox.Result{
				ID:     evs[1].GetId().Into(),
				Status: outbox.StatusPublished,
			},
		)
	})
	fmt.Println(err)

	// Output:
	// 2
	// true
	// true
	// <nil>
	// 2
	// true
	// true
	// <nil>
}

func Example_mockEventRepoAbortTransaction() {
	repo := newRepo()
	err := repo.Process(context.TODO(), func(eb outbox.EventsBatch) error {
		evs, err := eb.Poll(outbox.Arg{Limit: 2})
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(len(evs))
		fmt.Println(evs[0].GetId().Into() == events[0].GetId().Into())
		fmt.Println(evs[1].GetId().Into() == events[1].GetId().Into())

		return fmt.Errorf(`a dedicated error`)
	})
	fmt.Println(err)

	err = repo.Process(context.TODO(), func(eb outbox.EventsBatch) error {
		evs, err := eb.Poll(outbox.Arg{Limit: 2})
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(len(evs))
		fmt.Println(evs[0].GetId().Into() == events[0].GetId().Into())
		fmt.Println(evs[1].GetId().Into() == events[1].GetId().Into())

		return eb.AdvanceCursor(
			outbox.Result{
				ID:     evs[0].GetId().Into(),
				Status: outbox.StatusPublished,
			},
			outbox.Result{
				ID:     evs[1].GetId().Into(),
				Status: outbox.StatusPublished,
			},
		)
	})
	fmt.Println(err)

	// Output:
	// 2
	// true
	// true
	// a dedicated error
	// 2
	// true
	// true
	// <nil>
}
