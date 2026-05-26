package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*VolumeReordered)(nil)

type VolumeReordered struct {
	id, aid, VID ID
	at           time.Time
	Seq          int
}

func NewVolumeReordered(aid, vid ID, seq int) VolumeReordered {
	return VolumeReordered{
		id:  ulid.Make(),
		aid: aid,
		VID: vid,
		at:  time.Now(),
		Seq: seq,
	}
}

func (v VolumeReordered) AggID() ID       { return v.aid }
func (v VolumeReordered) AggKind() string { return `Novel` }
func (v VolumeReordered) EmittedAt() time.Time {
	panic("unimplemented")
}

func (v VolumeReordered) ID() ID {
	panic("unimplemented")
}

func (v VolumeReordered) Kind() string {
	panic("unimplemented")
}

func (v VolumeReordered) Payload() ([]byte, error) {
	panic("unimplemented")
}

func (v VolumeReordered) String() string {
	panic("unimplemented")
}
