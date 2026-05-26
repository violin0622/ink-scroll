package vo

import (
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Can be either:
// - Part-Chapter
// - Chapter
type TOC struct {
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

func (t *TOC) RenameVolume(id ID, title string) error {
	i := slices.IndexFunc(
		t.volumes,
		func(v Volume) bool { return v.ID == id })
	if i < 0 {
		return ErrVolumeNotFound
	}

	t.volumes[i].Title = title
	return nil
}

// ReorderVolume reset target volume to [seq].
// seq must in the range [0, len(volumes))
func (t TOC) ReorderVolume(id ID, seq int) (TOC, error) {
	if seq < 0 || seq >= len(t.volumes) {
		return t, ErrInvalidVolumeSequence
	}

	idx := t.volumeByID(id)
	if idx < 0 {
		return t, ErrVolumeNotFound
	}

	v := t.volumes[idx]
	switch {
	case idx > seq: // move forward
		for i := idx - 1; i >= seq; i-- {
			t.volumes[i+1] = t.volumes[i]
		}
	case idx < seq: // move backward
		copy(t.volumes[idx:], t.volumes[idx+1:seq+1])
	default:
	}
	t.volumes[seq] = v

	return t, nil
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

func (t TOC) ChaptersByVolume(id ID) []int {
	if i := t.volumeByID(id); i < 0 {
		return nil
	}

	chapters := []int{}
	for i := range t.chapters {
		if t.chapters[i].VolumeID == id {
			chapters = append(chapters, i)
		}
	}

	return chapters
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

	t.volumes[i].Description = desc
	return t, nil
}

// volume name must be unique.
func (t *TOC) NewVolume(
	id ID,
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
		ID:          id,
		Title:       title,
		Description: desc,
		Kind:        kind,
	})

	return nil
}

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

// func (t TOC) noRepeatedOrLostedChapterSequence() error {
// 	for seq := range len(t.chapters) {
// 		if t.chapters[seq].Sequence != seq {
// 			return status.New(codes.Internal, `chapter sequence missed or repeated`).
// 				Err()
// 		}
// 	}
//
// 	return nil
// }

// func (t TOC) chapterSequenceOrdered() error {
// 	if slices.IsSortedFunc(t.chapters, ChapterCmp) {
// 		return nil
// 	}
// 	return status.New(codes.Internal, `chapter sequence not ordered`).Err()
// }
