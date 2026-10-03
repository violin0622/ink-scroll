package db

import (
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/violin0622/ink-scroll/faculty/domain/teacher"
)

type teacher struct {
	Model[ulid.ULID]

	Name string
}

func (t *teacher) From(tc *domain.Teacher) {
	if t == nil {
		return
	}
	t.ID = tc.ID()
	t.Name = tc.Name()
}

var _ domain.TeacherRepo = (*teacherRepo)(nil)

type teacherRepo struct {
	db *gorm.DB
}

// Create implements domain.Repo.
func (tr *teacherRepo) Create(
	ctx context.Context, t *domain.Teacher,
) error {
	var mt teacher
	mt.From(t)

	return tr.db.Model(&teacher{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&mt).
		Error
}

// Delete implements domain.Repo.
func (t *teacherRepo) Delete(context.Context, ulid.ULID) error {
	panic("unimplemented")
}

// Get implements domain.Repo.
func (t *teacherRepo) Get(
	context.Context,
	ulid.ULID,
	...domain.GetOption,
) (*domain.Teacher, error) {
	panic("unimplemented")
}

// Update implements domain.Repo.
func (t *teacherRepo) Update(context.Context, *domain.Teacher) error {
	panic("unimplemented")
}
