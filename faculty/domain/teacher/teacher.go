package teacher

import (
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
	ev "github.com/violin0622/ink-scroll/faculty/domain/teacher/events"
	"github.com/violin0622/ink-scroll/faculty/domain/teacher/vo"
)

var _ Aggregate = (*Teacher)(nil)

type Teacher struct {
	id       ulid.ULID
	name     string
	gender   Gender
	birthday Birthday
	hs       HireStatus

	repo EventRepo
}

// 所有字段都是 private，仅暴露只读方法。
// 这样可以保证只有当前包内的代码可以修改实体的状态。
// 不要定义诸如 SetID, SetName 之类的方法，那样跟
// public 就没有区别了。
func (t Teacher) ID() ulid.ULID  { return t.id }
func (t Teacher) Kind() string   { return `Teacher` }
func (t Teacher) Name() string   { return t.name }
func (t Teacher) Gender() Gender { return t.gender }
func (t Teacher) Age() uint8 {
	age, _ := t.birthday.Age()
	return uint8(age)
}

func (t *Teacher) ReplayEvents(e ...ev.Event) error {
	if len(e) == 0 {
		return nil
	}

	if e[0].Kind() != (ev.TeacherHired{}).Kind() {
		return common.ErrInitialEvent
	}

	for i := range e {
		if err := t.applyEvent(e[i]); err != nil {
			return fmt.Errorf(`replay events: %w`, err)
		}
	}

	return nil
}

func (t *Teacher) applyEvent(e ev.Event) error {
	switch e := e.(type) {
	case ev.TeacherHired:
		return t.applyHired(e)
	case ev.TeacherRetired:
		return t.applyRetired(e)
	default:
		return common.ErrUnknownEventKind(e.Kind())
	}
}

func (t *Teacher) applyHired(e ev.TeacherHired) error {
	et, ok := e.Payload().(*Teacher)
	if !ok {
		return fmt.Errorf(`invalid event payload`)
	}

	t.name = et.name
	t.birthday = et.birthday
	t.gender = et.gender

	return nil
}

func (t *Teacher) applyRetired(e ev.TeacherRetired) error {
	if e.AggID() != t.id {
		return fmt.Errorf(`ID mismatch`)
	}

	t.hs = vo.Retired
	return nil
}
