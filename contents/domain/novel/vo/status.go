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

type ReleaseStatus int

const (
	_ ReleaseStatus = iota
	Draft
	Reviewing
	Published
	Banned
)

var relStatusString = map[ReleaseStatus]string{
	Draft:     `Draft`,
	Reviewing: `Reviewing`,
	Published: `Published`,
	Banned:    `Banned`,
}

func (s ReleaseStatus) String() string { return cmp.Or(relStatusString[s], unknownStr) }

func (s ReleaseStatus) PublishReviewSubmitted() (ReleaseStatus, error) {
	if s == Draft {
		return Reviewing, nil
	}

	return s, invalidConvertionError(s, Reviewing)
}

func (s ReleaseStatus) PublishReviewApproved() (ReleaseStatus, error) {
	if s == Reviewing {
		return Published, nil
	}
	return s, invalidConvertionError(s, Published)
}

func (s ReleaseStatus) PublishReviewRejected() (ReleaseStatus, error) {
	if s == Reviewing {
		return Draft, nil
	}
	return s, invalidConvertionError(s, Draft)
}

func (s ReleaseStatus) Banned() (ReleaseStatus, error) {
	if s == Published {
		return Banned, nil
	}
	return s, invalidConvertionError(s, Banned)
}

func (s ReleaseStatus) AppealApproved() (ReleaseStatus, error) {
	if s == Banned {
		return Published, nil
	}
	return s, invalidConvertionError(s, Published)
}
