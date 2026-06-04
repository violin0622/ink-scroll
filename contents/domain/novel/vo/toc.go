package vo

import (
	"fmt"
	"slices"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = fmt.Println

// 公有方法设计原则：
// 简单、正交。
type TOC struct {
	// chapters 按 VolumeID 排列，顺序与 volumes 顺序相同。
	chapters []Chapter
	volumes  []Volume
}

func NewTOC() (TOC, error) {
	toc := TOC{volumes: []Volume{newVolume()}}
	return toc.maybeOk()
}

func RestoreTOC(v []Volume, c []Chapter) (TOC, error) {
	toc := TOC{chapters: c, volumes: v}
	return toc.maybeOk()
}

func (t TOC) AppendChapter(c Chapter) (TOC, error) {
	vi := t.volumeByID(c.VolumeID)
	if vi < 0 {
		return TOC{}, ErrVolumeNotFound
	}

	_, r := t.chapterRangeByVolumeIndex(vi)

	t.chapters = append(t.chapters, Chapter{})
	rcopy(t.chapters[r:], t.chapters)
	t.chapters[r] = c

	return t.maybeOk()
}

func (t TOC) DropChapter(cid ID) (TOC, error) {
	i := t.chapterByID(cid)
	if i < 0 {
		return TOC{}, ErrChapterNotFound
	}

	copy(t.chapters[i:], t.chapters[i+1:])

	return t.maybeOk()
}

// ReorderChapter change chapter sequence in same volume (relative sequence).
// seq is relative sequence in the volume. A too large or small seq cause
// chapter move to last or first place in the volume, respectively.
func (t TOC) ReorderChapter(cid ID, seq int) (TOC, error) {
	i := t.chapterByID(cid)
	if i < 0 {
		return t.err(ErrChapterNotFound)
	}

	vid := t.volumeByID(t.chapters[i].VolumeID)
	l, r := t.chapterRangeByVolumeIndex(vid)
	if l < 0 {
		// In a normal toc, this branch should not reach.
		return t.err(ErrVolumeNotFound)
	}

	// convert relative seq to absolute seq.
	curSeq, targetSeq := i, max(min(l+seq, r-1), l)

	switch {
	case curSeq < targetSeq:
		c := t.chapters[curSeq]
		copy(t.chapters[curSeq:], t.chapters[curSeq+1:targetSeq+1])
		t.chapters[targetSeq] = c
	case targetSeq < curSeq:
		c := t.chapters[targetSeq]
		copy(t.chapters[targetSeq:], t.chapters[targetSeq+1:curSeq+1])
		t.chapters[curSeq] = c
	default:
		// no op
	}

	return t.maybeOk()
}

// MoveToVolume move chapter to target volume vid, this may cause sequence
// change.
// chapter will be set in last place in the target volume.
func (t TOC) MoveToVolume(cid, vid ID) (TOC, error) {
	ci := t.chapterByID(cid)
	if ci < 0 {
		return t.err(ErrChapterNotFound)
	}

	vi := t.volumeByID(vid)
	if vi < 0 {
		return t.err(ErrVolumeNotFound)
	}

	l, r := t.chapterRangeByVolumeIndex(vi)

	// already in vid Volume, noop
	if ci >= l && ci < r {
		return t, nil
	}

	switch {
	case ci < l: // move left to right.
		c := t.chapters[ci].ChangeVolume(vid)
		copy(t.chapters[ci:r], t.chapters[ci+1:r])
		t.chapters[r-1] = c
	case ci >= r: // move right to left.
		c := t.chapters[ci].ChangeVolume(vid)
		rcopy(t.chapters[r:ci+1], t.chapters[r-1:ci])
		t.chapters[max(l, r-1)] = c
	}

	return t.maybeOk()
}

func (t TOC) RenameVolume(vid ID, title string) (TOC, error) {
	i := t.volumeByID(vid)
	if i < 0 {
		return t, ErrVolumeNotFound
	}

	t.volumes[i] = t.volumes[i].Rename(title)
	return t, t.checkInvariants()
}

// ReorderVolume reset target volume to seq,
// all chapters in the volume will be reset too.
// seq must in the range [0, len(volumes))
func (t TOC) ReorderVolume(vid ID, targetSeq int) (TOC, error) {
	if targetSeq < 0 || targetSeq >= len(t.volumes) {
		return t.err(ErrInvalidVolumeSequence)
	}

	curSeq := t.volumeByID(vid)
	if curSeq < 0 {
		return t.err(ErrVolumeNotFound)
	}

	switch {
	case curSeq < targetSeq:
		t.moveLeftToRightV2(curSeq, targetSeq)
	case targetSeq < curSeq:
		t.moveRightToLeftV2(targetSeq, curSeq)
	default:
	}

	return t.maybeOk()
}

func (t *TOC) moveLeftToRightV2(l, r int) {
	ll, lr := t.chapterRangeByVolumeIndex(l)
	rl, rr := t.chapterRangeByVolumeIndex(r)

	v := t.volumes[l]
	copy(t.volumes[l:], t.volumes[l+1:r+1])
	t.volumes[r] = v

	if ll == lr {
		return
	}
	if rl == rr {
	}

	tmp := make([]Chapter, lr-ll)
	copy(tmp, t.chapters[ll:lr])
	copy(t.chapters[ll:], t.chapters[lr:rr])
	rcopy(t.chapters[:rr], tmp)
}

// chapterRangeByVolumeIndex returns a left-open and right-closed index
// interval.
// of given volumes at vi.
func (t *TOC) chapterRangeByVolumeIndex(vi int) (l, r int) {
	if vi < 0 || vi >= len(t.volumes) {
		return -1, -1
	}

	var v int
	var ok bool // whether we found a l index.

	for i := 0; i < len(t.chapters); {
		if ok {
			if v != vi {
				r = i
				break
			}
		} else {
			if v == vi {
				l = i
				ok = true
			}
		}

		if t.chapters[i].VolumeID != t.volumes[v].ID {
			v += 1
		} else {
			i += 1
		}
	}

	if v != vi {
		return l, r
	} else {
		return l, len(t.chapters)
	}
}

func (t *TOC) moveRightToLeftV2(l, r int) {
	ll, _ := t.chapterRangeByVolumeIndex(l)
	rl, rr := t.chapterRangeByVolumeIndex(r)

	v := t.volumes[r]
	rcopy(t.volumes[:r+1], t.volumes[l:r])
	t.volumes[l] = v

	if rr-rl == 0 {
		return
	}

	tmp := make([]Chapter, rr-rl)
	rcopy(tmp, t.chapters[rl:rr])
	rcopy(t.chapters[:rr], t.chapters[ll:rl])
	copy(t.chapters[ll:], tmp)
}

// rcopy is like copy in std, but from right to left.
func rcopy[T any](dst, src []T) {
	for i, j := len(dst)-1, len(src)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		dst[i] = src[j]
	}
}

// DropVolume deletes a volume.
// REQ-VOLUME-011：空的自定义分卷可以删除。
// REQ-VOLUME-012：存在已发布章节的分卷不能直接删除。
func (t TOC) DropVolume(vid ID) (TOC, error) {
	i := t.volumeByID(vid)
	if i < 0 {
		return t, ErrVolumeNotFound
	}

	if len(t.volumes) == 1 {
		return t, ErrDropLastVolume
	}

	chapters := t.ChaptersByVolume(vid)
	if len(chapters) != 0 {
		return t, VolumeHasChapterError{len(chapters)}
	}

	copy(t.volumes[i:], t.volumes[i+1:])
	t.volumes = t.volumes[:len(t.volumes)-1]

	return t, nil
}

// ChaptersByVolume returns chapters index in the vid Volume.
func (t TOC) ChaptersByVolume(vid ID) []int {
	if i := t.volumeByID(vid); i < 0 {
		return nil
	}

	return lo.FilterMap(t.chapters, func(c Chapter, i int) (int, bool) {
		if c.VolumeID == vid {
			return i, true
		}
		return 0, false
	})
}

func (t TOC) volumeByID(id ID) int {
	return slices.IndexFunc(
		t.volumes,
		func(v Volume) bool { return v.ID == id },
	)
}

func (t TOC) UpdateVolumeDesc(id ID, desc string) (TOC, error) {
	i := slices.IndexFunc(
		t.volumes,
		func(v Volume) bool { return v.ID == id })
	if i < 0 {
		return t, ErrVolumeNotFound
	}

	t.volumes[i] = t.volumes[i].UpdateDescription(desc)
	return t, nil
}

// NewVolume appends a new volume to TOC.
// volume name must be unique.
func (t TOC) AppendVolume(
	title, desc string,
	kind VolumeKind,
) (TOC, Volume, error) {
	i := slices.IndexFunc(t.volumes, func(v Volume) bool {
		return v.Title == title
	})
	if i >= 0 {
		return TOC{}, Volume{}, VolumeTitleAlreadyExistError{
			ID:    t.volumes[i].ID,
			Title: title,
		}
	}

	v := Volume{
		ID:          ulid.Make(),
		Title:       title,
		Description: desc,
		Kind:        kind,
	}
	t.volumes = append(t.volumes, v)

	return t, v, nil
}

// RestoreVolume are intend to be internal used when restore data from repo.
func (t *TOC) RestoreVolume(
	id ID,
	title, desc string,
	kind VolumeKind,
) {
	t.volumes = append(t.volumes, Volume{
		ID:          id,
		Title:       title,
		Description: desc,
		Kind:        kind,
	})
}

func (t TOC) LatestChapter() (c Chapter, err error) {
	if t.isEmpty() {
		return c, ErrNoContent
	}

	return t.chapters[len(t.chapters)-1], nil
}

func (t TOC) Volumes() (v []Volume) {
	v = make([]Volume, len(t.volumes))
	copy(v, t.volumes)
	return
}

func (t TOC) GetChapters() (c []Chapter) {
	if t.isEmpty() {
		return nil
	}

	c = make([]Chapter, len(t.chapters)-1)
	copy(c, t.chapters[1:])
	return c
}

func (t TOC) ChapterCount() int {
	return len(t.chapters) - 1
}

func (t TOC) ChapterSequenceInVolume(cid ID) int {
	i := t.chapterByID(cid)
	if i < 0 {
		return 0
	}

	n := 0
	vid := t.chapters[i].VolumeID
	for _, c := range t.chapters[:i+1] {
		if c.VolumeID == vid {
			n += 1
		}
	}

	return n
}

func (t TOC) chapterByID(cid ID) int {
	return slices.IndexFunc(
		t.chapters,
		func(c Chapter) bool { return c.ID == cid },
	)
}

// func (t TOC) isValidSequence(seq int) bool {
// 	return seq >= 1 && seq < len(t.chapters)
// }

func (t TOC) isEmpty() bool { return len(t.chapters) == 1 }

func (t TOC) checkInvariants() error {
	invariants := []func() error{
		// t.beginFromSentinel,
		t.atleastOneMainVolume,
		t.chaptersSortedByVolumeID,
	}

	for _, fn := range invariants {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// 新书首章必须在正文卷
// TBD:
// Case: 上传了两章，并将第一章移动至其他分卷？
// Case: 上传一章、删除一章、再上传一章？
func (t TOC) firstChapterGoesMainVolume() error {
	return nil
}

func (t TOC) atleastOneMainVolume() error {
	if ok := lo.ContainsBy(t.volumes, func(v Volume) bool {
		return v.Kind == VolumeKindMain
	}); !ok {
		return InvariantBrokenError{
			Msg: `MUST have at least one volume of main kind.`,
		}
	}
	return nil
}

func (t TOC) beginFromSentinel() error {
	if len(t.chapters) < 1 ||
		t.chapters[0] != SentinelChapter {
		return status.New(codes.Internal, `TOC structure destoyed`).Err()
	}
	return nil
}

// chaptersSortedByVolumeID 检查如下不变性：
// 1. A 中的所有元素都属于 B；
// 2. A 中的所有重复元素都相邻；
// 3. A 中的元素顺序与 B 相同。
//
// 例如，
//
// A=[1,1,1,2,2,4,4] B = [1,2,4]  -->合法
// A=[1,1,1,1,2,2,4,4] B=[1,2,3,4] --> 合法
// A=[1,1,1,1,4,4,2] B=[1,4,2] --> 合法
// A=[1,1,1,2,2,4,2] B=[1,2,4] --> 非法
// A=[0,1,1,1,2,2,4,4] B=[1,2,4] --> 非法
// A=[1,1,1,1,4,4,2,2] B=[1,2,4] --> 非法
func (t TOC) chaptersSortedByVolumeID() error {
	for c, v := 0, 0; c < len(t.chapters); {
		if t.chapters[c].VolumeID == t.volumes[v].ID {
			c += 1
			continue
		}

		v += 1

		if v >= len(t.volumes) {
			return InvariantBrokenError{
				Msg: `chapters MUST be arranged in the same way as volumes`,
			}
		}
	}

	return nil
}

func (t TOC) err(e error) (TOC, error) { return TOC{}, e }
func (t TOC) maybeOk() (TOC, error) {
	if e := t.checkInvariants(); e != nil {
		return t.err(e)
	}
	return t, nil
}
