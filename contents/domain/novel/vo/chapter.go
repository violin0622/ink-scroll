package vo

import (
	"github.com/oklog/ulid/v2"
)

var SentinelChapter = Chapter{}

type (
	ID      = ulid.ULID
	Chapter struct {
		ID, VolumeID ID

		// title are used for test/debug purpose(easy to read), do not use in
		// bussiness logic.
		title string
	}
)

func NewChapter(id, vid ID) Chapter {
	return Chapter{ID: id, VolumeID: vid}
}

func (c Chapter) ChangeVolume(vid ID) Chapter {
	c.VolumeID = vid
	return c
}
