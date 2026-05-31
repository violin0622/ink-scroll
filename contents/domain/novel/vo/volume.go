package vo

type Volume struct {
	ID ID

	Title       string
	Description string
	Kind        VolumeKind
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

const (
	_ VolumeKind = iota
	Main
	Extra
	Related
)

var volumeKindStr = map[VolumeKind]string{
	Main:    `Main`,
	Extra:   `Extra`,
	Related: `Related`,
}

func (k VolumeKind) String() string { return volumeKindStr[k] }
