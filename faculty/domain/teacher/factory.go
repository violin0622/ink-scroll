package teacher

import (
	"context"
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"

	ev "github.com/violin0622/ink-scroll/faculty/domain/teacher/events"
	"github.com/violin0622/ink-scroll/faculty/domain/teacher/vo"
)

type teacherBuilder struct {
	t   Teacher
	err error
}

type teacherFactory struct {
	er EventRepo
}

func NewTeacherFactory(er EventRepo) teacherFactory {
	return teacherFactory{er: er}
}

func (tf teacherFactory) NewHire(
	ctx context.Context, name, gender, birthday string) (
	t Teacher, err error,
) {
	t, err = tf.Builder(ulid.Make()).
		Name(name).
		Gender(gender).
		Birthday(birthday).
		Build()
	if err != nil {
		return t, err
	}

	t.hs = vo.NewHire

	e := ev.NewTeacherHired(t.ID(), t)
	if err = tf.er.Append(ctx, e); err != nil {
		return t, err
	}

	return t, err
}

func (tf teacherFactory) Builder(id ulid.ULID) *teacherBuilder {
	return &teacherBuilder{
		t: Teacher{
			id: id,

			repo: tf.er,
		},
	}
}

func (tb *teacherBuilder) Name(n string) *teacherBuilder {
	if tb == nil || tb.err != nil {
		return tb
	}

	if len(strings.Trim(n, ` `)) == 0 {
		tb.err = errors.New(`name can't be empty`)
	}

	tb.t.name = n
	return tb
}

func (tb *teacherBuilder) Gender(g string) *teacherBuilder {
	if tb == nil || tb.err != nil {
		return tb
	}

	tb.t.gender, tb.err = vo.ParseGender(g)
	return tb
}

func (tb *teacherBuilder) Birthday(b string) *teacherBuilder {
	if tb == nil || tb.err != nil {
		return tb
	}

	tb.t.birthday, tb.err = vo.ParseBirthday(b)
	return tb
}

func (tb *teacherBuilder) Build() (t Teacher, err error) {
	return tb.t, tb.err
}
