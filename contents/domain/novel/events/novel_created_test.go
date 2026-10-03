package events_test

import (
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ne "github.com/violin0622/ink-scroll/contents/domain/novel/events"
	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	eventv1 "github.com/violin0622/ink-scroll/proto/gen/go/event/v1"
	ulidpb "github.com/violin0622/ink-scroll/proto/gen/go/ulid"
)

var _ = Describe(`NovelCreatedV2 event fromPB`, func() {
	var (
		payload, _ = anypb.New(&novelv1.NovelCreatedPayload{
			AuthorId:    ulidpb.From(ulid.Make()),
			Title:       `My Book Title`,
			Description: `My book is good!`,
			Tags:        []string{`Tokyo`, `Hot`},
			Category:    `Story`,
		})
		pbEvent = eventv1.Event{
			Id:            ulidpb.From(ulid.Make()),
			AggregateId:   ulidpb.From(ulid.Make()),
			EmitAt:        timestamppb.Now(),
			Kind:          "created",
			AggregateKind: "novel",
			Payload:       payload,
		}
	)

	It(`can unmarshal from proto`, func() {
		var e ne.NovelCreated
		Ω(e.FromPB(&pbEvent)).Should(Succeed())
		Ω(e.ID()).Should(BeEquivalentTo(pbEvent.GetId().Into()))

		raw, _ := protojson.Marshal(pbEvent.GetPayload())
		Ω(e.Payload()).Should(Equal(raw))
	})
})
