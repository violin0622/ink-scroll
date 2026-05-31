package vo

import (
	"errors"
	"fmt"
	"slices"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Can be either:
// - Part-Chapter
// - Chapter
type TOC struct {
	// chapters 按 VolumeID 排列，顺序与 volumes 顺序相同。
	chapters []Chapter
	volumes  []Volume
}

func NewTOC(c ...Chapter) (TOC, error) {
	toc := TOC{
		chapters: append([]Chapter{SentinelChapter}, c...),
		volumes: []Volume{{
			ID:          ID{},
			Title:       "",
			Description: "",
			Kind:        0,
		}},
	}
	return toc, toc.checkInvariants()
}

func (t TOC) AppendChapter() { panic(`unimplement`) }
func (t TOC) DropChapter()   { panic(`unimplement`) }

// ReorderChapter change chapter sequence in same volume (relative sequence).
func (t TOC) ReorderChapter() { panic(`unimplement`) }

// MoveToVolume move chapter to target volume, this may cause sequance change.
func (t TOC) MoveToVolume(cid, vid ID) { panic(`unimplement`) }

func (t TOC) RenameVolume(vid ID, title string) (TOC, error) {
	i := t.volumeByID(vid)
	if i < 0 {
		return t, ErrVolumeNotFound
	}

	t.volumes[i] = t.volumes[i].Rename(title)
	return t, nil
}

// ReorderVolume reset target volume to seq,
// all chapters in the volume will be reset too.
// seq must in the range [0, len(volumes))
func (t TOC) ReorderVolume(vid ID, targetSeq int) (TOC, error) {
	if targetSeq < 0 || targetSeq >= len(t.volumes) {
		return t, ErrInvalidVolumeSequence
	}

	curSeq := t.volumeByID(vid)
	if curSeq < 0 {
		return t, ErrVolumeNotFound
	}

	fmt.Println(t.volumes)
	fmt.Println(t.chapters)

	switch {
	case curSeq < targetSeq:
		t.moveLeftToRight(curSeq, targetSeq)
	case targetSeq < curSeq:
		t.moveLeftToRight(targetSeq, curSeq)
	default:
	}

	return t, nil
}

func (t *TOC) moveLeftToRight(l, r int) {
	v := t.volumes[l]
	copy(t.volumes[l:], t.volumes[l+1:r+1])
	t.volumes[r] = v

	lChpIdxs := t.ChaptersByVolume(t.volumes[l].ID)
	rChpIdxs := t.ChaptersByVolume(t.volumes[r].ID)
	ll, lr := lChpIdxs[0], lChpIdxs[len(lChpIdxs)-1]
	rl, rr := rChpIdxs[0], rChpIdxs[len(rChpIdxs)-1]

	fmt.Println(l, r)
	fmt.Println(lChpIdxs, rChpIdxs)
	fmt.Println(ll, lr, rl, rr)
	tmp := make([]Chapter, len(lChpIdxs))
	copy(tmp, t.chapters[ll:lr+1])
	copy(t.chapters[ll:], t.chapters[lr+1:rr+1])
	rcopy(t.chapters[:rr+1], tmp)
}

func (t *TOC) moveRightToLeft(l, r int) {
	v := t.volumes[r]
	rcopy(t.volumes[:r+1], t.volumes[l:r])
	t.volumes[l] = v

	lChpIdxs := t.ChaptersByVolume(t.volumes[l].ID)
	rChpIdxs := t.ChaptersByVolume(t.volumes[r].ID)
	ll, _ := lChpIdxs[0], lChpIdxs[len(lChpIdxs)-1]
	rl, rr := rChpIdxs[0], rChpIdxs[len(rChpIdxs)-1]

	tmp := make([]Chapter, len(rChpIdxs))
	rcopy(tmp, t.chapters[rl:rr+1])
	rcopy(t.chapters[:rr+1], t.chapters[ll:rl+1])
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

func (t TOC) ChaptersByVolume(vid ID) []int {
	if i := t.volumeByID(vid); i < 0 {
		return nil
	}

	chp := lo.Filter(
		t.chapters,
		func(c Chapter, _ int) bool { return c.VolumeID == vid })
	return lo.Map(
		chp,
		func(_ Chapter, i int) int { return i })
}

func (t TOC) volumeByID(id ID) int {
	return slices.IndexFunc(t.volumes, func(v Volume) bool { return v.ID == id })
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
func (t *TOC) NewVolume(
	title, desc string,
	kind VolumeKind,
) error {
	i := slices.IndexFunc(t.volumes, func(v Volume) bool {
		return v.Title == title
	})
	if i >= 0 {
		return VolumeTitleAlreadyExistError{
			ID:    t.volumes[i].ID,
			Title: title,
		}
	}

	t.volumes = append(t.volumes, Volume{
		ID:          ulid.Make(),
		Title:       title,
		Description: desc,
		Kind:        kind,
	})

	return nil
}

// RestoreVolume are intend to be used when restore data from repo.
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

func (t TOC) Get(seq int) (c Chapter, err error) {
	if t.isValidSequence(seq) {
		return c, ErrInvalidChapterSequence
	}

	return t.chapters[seq], nil
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

// func (t TOC) NextChapter() (c Chapter) {
// 	c.Sequence = len(t.chapters)
//
// 	return c
// }

func (t TOC) Append(c Chapter) (TOC, error) {
	t.chapters = append(t.chapters, c)
	return t, t.checkInvariants()
}

func (t TOC) Pop() (nt TOC, c Chapter, err error) {
	if t.isEmpty() {
		return t, c, ErrNoContent
	}

	last := len(t.chapters) - 1
	c, t.chapters = t.chapters[last], t.chapters[:last]

	return t, c, t.checkInvariants()
}

func (t TOC) ChapterSequence(cid ID) int { return len(t.chapters) }
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
	return slices.IndexFunc(t.chapters, func(c Chapter) bool { return c.ID == cid })
}

// Revise applies title and wordCount to chapter at seq sequence,
// returns original chapter.
// func (t TOC) Revise(c Chapter) (nt TOC, origin Chapter, err error) {
// 	if !t.isValidSequence(c.Sequence) {
// 		return t, c, ErrInvalidChapterSequence
// 	}
//
// 	origin, t.chapters[c.Sequence] = t.chapters[c.Sequence], c
//
// 	return t, origin, t.checkInvariants()
// }

func (t TOC) isValidSequence(seq int) bool {
	return seq >= 1 && seq < len(t.chapters)
}

func (t TOC) isEmpty() bool { return len(t.chapters) == 1 }

func (t TOC) checkInvariants() error {
	invariants := []func() error{
		t.beginFromSentinel,
		t.atleastOneVolume,
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

func (t TOC) atleastOneVolume() error {
	if len(t.chapters) < 2 {
		return status.New(codes.Internal, `at least one chapter`).Err()
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

func (t TOC) chaptersInVolume() error {
	const InvariantNum = `I`

	ok := slices.ContainsFunc(
		t.chapters,
		func(c Chapter) bool { return t.volumeByID(c.VolumeID) == -1 })

	if !ok {
		return InvariantBrokenError{
			Msg: "Chapter has invalid VolumeID.",
			Num: "InvariantNum",
		}
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
		// fmt.Println(c, t.chapters[c].VolumeID, v, t.volumes[v].ID)

		if t.chapters[c].VolumeID == t.volumes[v].ID {
			c += 1
			continue
		}

		v += 1

		if v >= len(t.volumes) {
			return errors.New(`invalid`)
		}
	}

	return nil
}
