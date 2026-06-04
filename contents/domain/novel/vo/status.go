package vo

import (
	"cmp"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ fmt.Stringer = (Status)(0)

type Status int

const (
	_ Status = iota
	Serial
	Desolated
	Completed
)

var statusString = map[Status]string{
	Serial:    `Serial`,
	Desolated: `Desolated`,
	Completed: `Completed`,
}

func (s Status) String() string { return cmp.Or(statusString[s], unknownStr) }

func (s Status) ChapterPublished() (Status, error) {
	switch s {
	case Serial, Desolated:
		return Serial, nil
	case Completed:
		return s, status.Error(
			codes.FailedPrecondition,
			`can't publish new chapter to a completed novel.`,
		)
	default:
		return s, status.Errorf(
			codes.FailedPrecondition,
			`unknown status '%s'`,
			s,
		)
	}
}

func (s Status) MarkCompleted() (Status, error) {
	return Completed, nil
}

func (s Status) MarkDesolated() (Status, error) {
	return Desolated, nil
}

type ReleaseStatus string

const (
	_         ReleaseStatus = ``
	Draft     ReleaseStatus = `Draft`
	Reviewing ReleaseStatus = `Reviewing`
	Released  ReleaseStatus = `Released`
	Banned    ReleaseStatus = `Banned`
)

func (s ReleaseStatus) NovelReviewSubmitted() (ReleaseStatus, error) {
	if s == Draft {
		return Reviewing, nil
	}

	return s, invalidConvertionError(s, Reviewing)
}

func (s ReleaseStatus) NovelReviewApproved() (ReleaseStatus, error) {
	if s == Reviewing {
		return Released, nil
	}
	return s, invalidConvertionError(s, Released)
}

func (s ReleaseStatus) NovelReviewRejected() (ReleaseStatus, error) {
	if s == Reviewing {
		return Draft, nil
	}
	return s, invalidConvertionError(s, Draft)
}

func (s ReleaseStatus) Banned() (ReleaseStatus, error) {
	if s == Released {
		return Banned, nil
	}
	return s, invalidConvertionError(s, Banned)
}

func (s ReleaseStatus) AppealApproved() (ReleaseStatus, error) {
	if s == Banned {
		return Released, nil
	}
	return s, invalidConvertionError(s, Released)
}

// var (
// 	ErrMutateCompletedNovel = status.New(
// 		codes.FailedPrecondition,
// 		`mutate a completed novel is not allowed`).Err()
// )
