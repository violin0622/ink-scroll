package draft

import (
	"time"

	"github.com/oklog/ulid/v2"
)

func NewDraft(title, main, ex string) (*Draft, error) {
	return &Draft{
		id:           ulid.Make(),
		title:        title,
		mainContent:  main,
		extraContent: ex,
		releaseAt:    time.Now(),
	}, nil
}
