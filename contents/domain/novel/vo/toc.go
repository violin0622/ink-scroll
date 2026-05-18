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
	toc := TOC{chapters: append([]Chapter{SentinelChapter}, c...)}
	return toc, toc.checkInvariants()
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

func (t TOC) NextChapter() (c Chapter) {
	c.Sequence = len(t.chapters)

	return c
}

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

// Revise applies title and wordCount to chapter at seq sequence,
// returns original chapter.
func (t TOC) Revise(c Chapter) (nt TOC, origin Chapter, err error) {
	if !t.isValidSequence(c.Sequence) {
		return t, c, ErrInvalidChapterSequence
	}

	origin, t.chapters[c.Sequence] = t.chapters[c.Sequence], c

	return t, origin, t.checkInvariants()
}

func (t TOC) isValidSequence(seq int) bool {
	return seq >= 1 && seq < len(t.chapters)
}

func (t TOC) isEmpty() bool { return len(t.chapters) == 1 }

func (t TOC) checkInvariants() error {
	invariants := []func() error{
		t.beginFromSentinel,
		// t.atleastOneChapter,
		t.chapterSequenceOrdered,
		t.noRepeatedOrLostedChapterSequence,
	}

	for _, fn := range invariants {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// func (t TOC) atleastOneChapter() error {
// 	if len(t.chapters) < 2 {
// 		return status.New(codes.Internal, `at least one chapter`).Err()
// 	}
// 	return nil
// }

func (t TOC) beginFromSentinel() error {
	if len(t.chapters) < 1 ||
		t.chapters[0] != SentinelChapter {
		return status.New(codes.Internal, `TOC structure destoyed`).Err()
	}
	return nil
}

func (t TOC) noRepeatedOrLostedChapterSequence() error {
	for seq := range len(t.chapters) {
		if t.chapters[seq].Sequence != seq {
			return status.New(codes.Internal, `chapter sequence missed or repeated`).
				Err()
		}
	}

	return nil
}

func (t TOC) chapterSequenceOrdered() error {
	if slices.IsSortedFunc(t.chapters, ChapterCmp) {
		return nil
	}
	return status.New(codes.Internal, `chapter sequence not ordered`).Err()
}
