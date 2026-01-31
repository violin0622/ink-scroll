package event

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"play-ddd/common"
	"play-ddd/contents/domain/novel"
	"play-ddd/contents/infra/repository/pg/datatypes/ts"
	dtulid "play-ddd/contents/infra/repository/pg/datatypes/ulid"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

type (
	ID    = dtulid.ULID
	Event struct {
		ID            ID `gorm:"primarykey"`
		CreatedTS     ts.Timestamp
		UpdatedTS     ts.Timestamp
		AggregateID   ID
		AggregateKind string
		Kind          string
		Payload       []byte
		AggregateSeq  uint64
		Status        *string
		Reason        *string
		Seq           *uint64
	}
)

func (e *Event) FromDomain(de common.Event[novel.ID, novel.ID]) (err error) {
	if e == nil {
		return nil
	}

	e.ID = ID(de.ID())
	e.AggregateID = ID(de.AggID())
	e.Kind = de.Kind()
	e.AggregateKind = de.AggKind()
	e.CreatedTS = ts.From(de.EmittedAt())
	e.Payload, err = de.Payload()
	return err
}

func (e *Event) IntoPB() (pe *novelv1.Event, err error) {
	pe = &novelv1.Event{Payload: &anypb.Any{}}
	if err = protojson.Unmarshal(e.Payload, pe.GetPayload()); err != nil {
		return pe, err
	}

	pe.Id = ulidpb.From(e.ID.Into())
	pe.AggregateId = ulidpb.From(e.AggregateID.Into())
	pe.EmitAt = timestamppb.New(e.CreatedTS.Into())
	pe.Kind = e.Kind
	pe.AggregateKind = e.AggregateKind

	return pe, err
}
