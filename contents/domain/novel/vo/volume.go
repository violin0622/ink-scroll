package vo

import (
	"cmp"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type Volume struct {
	ID ID

	Title       string
	Description string
	Kind        VolumeKind
}

func newVolume() Volume {
	return Volume{
		ID:          ulid.Make(),
		Title:       `Main`,
		Description: ``,
		Kind:        VolumeKindMain,
	}
}

func (v Volume) Rename(t string) Volume {
	v.Description = t
	return v
}

func (v Volume) UpdateDescription(d string) Volume {
	v.Description = d
	return v
}

type VolumeKind int

var _ fmt.Stringer = (VolumeKind)(0)

const (
	_ VolumeKind = iota
	VolumeKindMain
	VolumeKindExtra
	VolumeKindRelated
)

var volumeKindStr = map[VolumeKind]string{
	VolumeKindMain:    `Main`,
	VolumeKindExtra:   `Extra`,
	VolumeKindRelated: `Related`,
}

func (k VolumeKind) String() string {
	return cmp.Or(volumeKindStr[k], unknownStr)
}
