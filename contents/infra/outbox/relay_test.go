package outbox_test

import (
	"context"
	"fmt"
	"testing"

	"play-ddd/contents/infra/outbox"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

const raw = `{"@type":"type.googleapis.com/contents.novel.v1.NovelCreated","authorId":{"ulid":"01KD2Q3EVM65QRM1M28S0ZM7C6"},"title":"My Book Title","description":"My book is good!","tags":["Tokyo","Hot"],"category":"Story"}`

func TestUnmarshal(t *testing.T) {
	var aa anypb.Any
	err := protojson.Unmarshal([]byte(raw), &aa)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(&aa)
	t.Fail()
}

func TestStartIdeomepont(t *testing.T) {
	// outbox.NewRelay(bus EventBus, repo EventRepo, opts ...Option)
}

var _ outbox.EventBus = mockEventBus{}

type mockEventBus struct{}

// Pub implements [outbox.EventBus].
func (mockEventBus) Pub(context.Context, []*novelv1.Event) error {
	return nil
}

var _ outbox.EventRepo = mockEventRepo{}

type mockEventRepo struct{}

// Process implements [outbox.EventRepo].
func (mockEventRepo) Process(
	_ context.Context,
	fn func(outbox.EventsBatch) error,
) error {
	return fn(&mockEventbatch{})
}

var _ outbox.EventsBatch = &mockEventbatch{}

type mockEventbatch struct {
	cursor int
}

// AdvanceCursor implements [outbox.EventsBatch].
func (eb *mockEventbatch) AdvanceCursor(rs ...outbox.Result) error {
	eb.cursor += len(rs)
	return nil
}

// PollEvents implements [outbox.EventsBatch].
func (eb *mockEventbatch) PollEvents(outbox.Arg) ([]outbox.Event, error) {
	panic("unimplemented")
}

type mockevent struct{}
