package vo

import (
	"strconv"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ustatus "play-ddd/utils/grpc/status"
)

const errInfoDomain = `contents.novel.volume`

var (
	ErrNoContent              = status.New(codes.OutOfRange, `no content`).Err()
	ErrInvalidChapterSequence = status.New(codes.OutOfRange, `chapter sequence out of range`).
					Err()
	ErrInvalidVolumeSequence = status.New(codes.OutOfRange, `chapter sequence out of range`).
					Err()
	ErrVolumeNotFound = status.New(codes.NotFound, `volume not found`).
				Err()
	ErrChapterNotFound = status.New(codes.NotFound, `chapter not found`).
				Err()
)

func invalidConvertionError(a, b any) error {
	return status.Errorf(
		codes.FailedPrecondition,
		`invalid status convertion '%s'->'%s'`, a, b)
}

func volumeNotFoundError(nid, vid ID) error {
	st := status.New(codes.NotFound, `Unable to drop volume.`)
	st, err := st.WithDetails(
		&epb.ErrorInfo{
			Reason: novelv1.NovelError_NOVEL_ERROR_VOLUME_NOT_DROPPABLE.String(),
			Domain: errInfoDomain,
			Metadata: map[string]string{
				`novel_id`:  nid.String(),
				`volume_id`: vid.String(),
			},
		})
	if err != nil {
		return err
	}

	return st.Err()
}

func volumeHasChapterError(nid, vid ID, chapterCnt int) error {
	st := status.New(codes.FailedPrecondition, `Unable to drop volume.`)
	st, err := st.WithDetails(
		&epb.ErrorInfo{
			Reason: novelv1.NovelError_NOVEL_ERROR_VOLUME_NOT_DROPPABLE.String(),
			Domain: errInfoDomain,
			Metadata: map[string]string{
				`novel_id`:  nid.String(),
				`volume_id`: vid.String(),
				`chapters`:  strconv.Itoa(chapterCnt),
			},
		},
		&epb.PreconditionFailure{
			Violations: []*epb.PreconditionFailure_Violation{
				{
					Type:        novelv1.NovelViolationType_NOVEL_VIOLATION_TYPE_VOLUME_CHAPTER.String(),
					Subject:     "/volumes/" + vid.String(),
					Description: "There're still chapters in the volume, please drop or move them first.",
				},
			},
		})
	if err != nil {
		return err
	}

	return st.Err()
}

func dropLastVolumeError(nid, vid ID) error {
	return ustatus.FailedPreconditionWithDetails(
		`Unable to drop volume.`,
		&epb.ErrorInfo{
			Reason: novelv1.NovelError_NOVEL_ERROR_VOLUME_NOT_DROPPABLE.String(),
			Domain: errInfoDomain,
			Metadata: map[string]string{
				`novel_id`:  nid.String(),
				`volume_id`: vid.String(),
			},
		},
		&epb.PreconditionFailure{
			Violations: []*epb.PreconditionFailure_Violation{{
				Type:        novelv1.NovelViolationType_NOVEL_VIOLATION_TYPE_VOLUME_EMPTY.String(),
				Subject:     "/volumes/" + nid.String(),
				Description: "There're still chapters in the volume, please drop or move them first.",
			}},
		})
}
