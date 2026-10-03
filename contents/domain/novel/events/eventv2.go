package events

import (
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	"github.com/violin0622/ink-scroll/utils/xerr"
)

type m[A any] interface {
	~*A

	proto.Message
}

type event[A any, Ap m[A]] struct {
	id, aid       ID
	at            time.Time
	kind, aggKind string
	payload       Ap
	payload2      *A
}

func (e *event[A, B]) FromPB(pe *novelv1.Event) (err error) {
	defer xerr.Expect(&err, `fromPB`)

	e.id = pe.GetId().Into()
	e.aid = pe.GetAggregateId().Into()
	e.at = pe.GetEmitAt().AsTime()
	e.kind = pe.GetKind()
	e.aggKind = pe.GetAggregateKind()
	e.payload = new(A)
	return pe.GetPayload().UnmarshalTo(e.payload)
}

func (t event[A, Ap]) AggID() ID            { return t.aid }
func (t event[A, Ap]) EmittedAt() time.Time { return t.at }
func (t event[A, Ap]) ID() ID               { return t.id }
func (t event[A, Ap]) Kind() string         { return `created` }
func (t event[A, Ap]) AggKind() string      { return `novel` }
func (t event[A, Ap]) String() string       { return formatEvent(t) }

// Payload wraps payload into an Any and marshal it to json.
// Wrapped by Any so that there is type hint in marshaled byte.
func (t event[A, Ap]) Payload() ([]byte, error) {
	a, err := anypb.New(t.payload)
	if err != nil {
		return nil, err
	}

	return protojson.Marshal(a)
}
