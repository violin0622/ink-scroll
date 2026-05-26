package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = VolumeDropped{}

type VolumeDropped struct {
	id, aid, VID ID
	at           time.Time
}

func NewVolumeDropped(aid, vid ID) VolumeDropped {
	return VolumeDropped{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
		VID: vid,
	}
}

// AggID implements [common.Event].
func (v VolumeDropped) AggID() ID {
	panic("unimplemented")
}

// AggKind implements [common.Event].
func (v VolumeDropped) AggKind() string {
	panic("unimplemented")
}

// EmittedAt implements [common.Event].
func (v VolumeDropped) EmittedAt() time.Time {
	panic("unimplemented")
}

// ID implements [common.Event].
func (v VolumeDropped) ID() ID {
	panic("unimplemented")
}

// Kind implements [common.Event].
func (v VolumeDropped) Kind() string {
	panic("unimplemented")
}

// Payload implements [common.Event].
func (v VolumeDropped) Payload() ([]byte, error) {
	panic("unimplemented")
}

// String implements [common.Event].
func (v VolumeDropped) String() string {
	panic("unimplemented")
}
