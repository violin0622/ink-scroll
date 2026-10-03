package grpc

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	ustatus "github.com/violin0622/ink-scroll/utils/grpc/status"
)

func toStatus(ep *error) {
	if ep == nil || *ep == nil {
		return
	}

	if ve, ok := errors.AsType[novel.VolumeError](*ep); ok {
		*ep = fromVolumeError(ve)
	}

	*ep = unknownErrorWithDebugStacktrace(*ep, 1)
}

const (
	novelVolumeErrInfoDomain = `contents.novel.volume`
)

func unknownErrorWithDebugStacktrace(e error, skip int) error {
	if e == nil {
		return nil
	}

	const maxStackDepth = 32
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip+2, pcs)

	frames := runtime.CallersFrames(pcs[:n])
	stackEntries := make([]string, 0, n)
	for {
		frame, more := frames.Next()
		stackEntries = append(
			stackEntries,
			fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line),
		)
		if !more {
			break
		}
	}

	return ustatus.Unknown(`Unknown error.`).
		WithDebugInfo(stackEntries, e.Error()).
		Err()
}

func fromVolumeError(ve novel.VolumeError) error {
	var vhcErr vo.VolumeHasChapterError
	var vtaeErr vo.VolumeTitleAlreadyExistError
	switch {
	case errors.Is(ve.Cause, vo.ErrDropLastVolume):
		return dropLastVolumeError(ve.NovelID.String(), ve.VolumeID.String())
	case errors.Is(ve.Cause, vo.ErrVolumeNotFound):
		return volumeNotFoundError(ve.NovelID.String(), ve.VolumeID.String())
	case errors.As(ve.Cause, &vhcErr):
		return volumeHasChapterError(
			ve.NovelID.String(),
			ve.VolumeID.String(),
			strconv.Itoa(vhcErr.Cnt))
	case errors.As(ve.Cause, &vtaeErr):
		return volumeTitleAlreadyExistError(
			ve.NovelID.String(),
			vtaeErr.ID.String(),
			vtaeErr.Title)
	default:
		return unknownErrorWithDebugStacktrace(ve, 2)
	}
}

func volumeNotFoundError(nid, vid string) error {
	return ustatus.NotFound(`Novel volume not found.`).
		WithErrorInfo(
			novelv1.ErrorReason_VOLUME_NOT_DROPPABLE.String(),
			novelVolumeErrInfoDomain,
			map[string]string{
				`novel_id`:  nid,
				`volume_id`: vid,
			}).
		Err()
}

func dropLastVolumeError(nid, vid string) error {
	return ustatus.FailedPrecondition(`Unable to drop volume.`).
		WithErrorInfo(
			novelv1.ErrorReason_VOLUME_NOT_DROPPABLE.String(),
			novelVolumeErrInfoDomain,
			map[string]string{
				`novel_id`:  nid,
				`volume_id`: vid,
			}).
		WithPreconditionViolation(
			novelv1.NovelViolationType_NOVEL_VIOLATION_TYPE_VOLUME_EMPTY.String(),
			"/volumes/"+nid,
			"There're still chapters in the volume, please drop or move them first.").
		Err()
}

func volumeHasChapterError(nid, vid, cnt string) error {
	return ustatus.FailedPrecondition(`Unable to drop volume.`).
		WithErrorInfo(
			novelv1.ErrorReason_VOLUME_NOT_DROPPABLE.String(),
			novelVolumeErrInfoDomain,
			map[string]string{
				`novel_id`:     nid,
				`volume_id`:    vid,
				`chapters_cnt`: cnt,
			}).
		WithPreconditionViolation(
			novelv1.NovelViolationType_NOVEL_VIOLATION_TYPE_VOLUME_CHAPTER.String(),
			"/volumes/"+vid,
			"There're still chapters in the volume, please drop or move them first.",
		).
		Err()
}

func volumeTitleAlreadyExistError(nid, vid, title string) error {
	return ustatus.AlreadyExists(`Volume title already exists.`).
		WithErrorInfo(
			novelv1.ErrorReason_TITLE_ALREADY_EXIST.String(),
			novelVolumeErrInfoDomain,
			map[string]string{
				`novel_id`:          nid,
				`another_volume_id`: vid,
				`title`:             title,
			}).
		WithResourceInfo("volume", vid, ``, ``).
		Err()
}
