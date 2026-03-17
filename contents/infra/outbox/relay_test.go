package outbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"play-ddd/contents/infra/outbox"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

var _ ulidpb.ULID

func TestStartIdeomepont(t *testing.T) {
	eb := mockEventBus{}
	er := newRepo()
	relay := outbox.NewRelay(eb, er)
	if err := relay.Start(); err != nil {
		t.Errorf(`first start should not error: %s`, err)
	}

	if !relay.Ready() {
		t.Error(`should be ready aftert started.`)
	}

	if err := relay.Start(); err != nil {
		t.Errorf(`first start should not error: %s`, err)
	}
}

var _ outbox.EventBus = mockEventBus{}

type mockEventBus struct{}

// Pub implements [outbox.EventBus].
func (mockEventBus) Pub(context.Context, []*novelv1.Event) error {
	return nil
}

var _ outbox.EventRepo = mockEventRepo{}

type mockEventRepo struct {
	ev []*novelv1.Event
}

func newRepo() mockEventRepo {
	return mockEventRepo{ev: events}
}

// Process implements [outbox.EventRepo].
func (r mockEventRepo) Process(
	ctx context.Context,
	fn func(outbox.EventsBatch) error,
) error {
	return fn(&mockEventbatch{r: r})
}

// Process implements [outbox.EventRepo].

var _ outbox.EventsBatch = &mockEventbatch{}

type mockEventbatch struct {
	r      mockEventRepo
	cursor int
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
