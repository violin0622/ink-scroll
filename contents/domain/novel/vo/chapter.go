package vo

import (
	"github.com/oklog/ulid/v2"
)

var SentinelChapter = Chapter{}

type (
	ID      = ulid.ULID
	Chapter struct {
		ID ID
		// Sequence int
		VolumeID ID

		Title string
		// UpdatedAt  time.Time
	}
)

func NewChapter(id, vid ID) Chapter {
	return Chapter{ID: id, VolumeID: vid}
}

func (c Chapter) ChangeVolume(vid ID) Chapter {
	c.VolumeID = vid
	return c
}

// func ChapterCmp(a, b Chapter) int {
// 	return cmp.Compare(a.Sequence, b.Sequence)
// }
