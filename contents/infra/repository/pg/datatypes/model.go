package datatypes

import "github.com/violin0622/ink-scroll/contents/infra/repository/pg/datatypes/ts"

type Model[A comparable] struct {
	ID        A `gorm:"primarykey"`
	CreatedTs ts.Timestamp
	UpdatedTs ts.Timestamp
	DeletedTs ts.Timestamp
}
