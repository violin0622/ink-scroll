package vo

type Volume struct {
	ID ID

	Title       string
	Description string
	Kind        VolumeKind
}

type VolumeKind int

const (
	_ VolumeKind = iota
	Main
	Extra
	Related
)
