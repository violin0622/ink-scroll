package vo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"

	"github.com/violin0622/ink-scroll/utils/xtest"
)

// setupToc 解析以 ABBCDDDFF 形式的字符串，返回对应结构的 TOC。
// 不存在的字母视为卷存在而章节不存在。
func setupToc(pattern string) (t TOC) {
	volTitle, chpTitle := 'A', 'a'
	t.volumes = []Volume{{ID: ulid.Make(), Title: `A`, Kind: VolumeKindMain}}

	for i := 0; i < len(pattern); {
		if pattern[i] != byte(volTitle) {
			volTitle += 1
			t.volumes = append(t.volumes, Volume{
				ID:    ulid.Make(),
				Title: string(volTitle),
			})
			continue
		}

		t.chapters = append(t.chapters, Chapter{
			ID:       ulid.Make(),
			VolumeID: t.volumes[len(t.volumes)-1].ID,
			title:    string(chpTitle),
		})
		chpTitle += 1
		i += 1
	}

	return t
}

// func Test_setupToc(t *testing.T) {
// 	for _, p := range []string{
// 		`A`,
// 		`AAA`,
// 		`ABB`,
// 		`B`,
// 		`BCD`,
// 		`BD`,
// 		`BDD`,
// 		`AAABCCDDDD`,
// 	} {
// 		fmt.Println("")
// 		fmt.Println(p)
// 		toc := setupToc(p)
// 		for i := range toc.volumes {
// 			fmt.Println(toc.volumes[i])
// 		}
// 		for i := range toc.chapters {
// 			fmt.Println(toc.chapters[i])
// 		}
// 	}
// }

func TestTOCReorderVolume(t *testing.T) {
	tests := []struct {
		name      string
		toc       TOC
		from, to  int
		expectVol string
		expectChp string
	}{{
		name:      `noop-0`,
		toc:       setupToc(`ABCDE`),
		from:      0,
		to:        0,
		expectVol: `ABCDE`,
		expectChp: `abcde`,
	}, {
		name:      `noop-2`,
		toc:       setupToc(`ABCDE`),
		from:      2,
		to:        2,
		expectVol: `ABCDE`,
		expectChp: `abcde`,
	}, {
		name:      `noop-4`,
		toc:       setupToc(`ABCDE`),
		from:      4,
		to:        4,
		expectVol: `ABCDE`,
		expectChp: `abcde`,
	}, {
		name:      `noop-ma`,
		toc:       setupToc(`AAABCDE`),
		from:      4,
		to:        4,
		expectVol: `ABCDE`,
		expectChp: `abcdefg`,
	}, {
		name:      `noop-m`,
		toc:       setupToc(`AAABCDE`),
		from:      4,
		to:        4,
		expectVol: `ABCDE`,
		expectChp: `abcdefg`,
	}, {
		name:      `noop-m`,
		toc:       setupToc(`BCDEEEE`),
		from:      4,
		to:        4,
		expectVol: `ABCDE`,
		expectChp: `abcdefg`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`ABCDE`),
		from:      0,
		to:        1,
		expectVol: `BACDE`,
		expectChp: `bacde`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`AAABCDE`),
		from:      0,
		to:        1,
		expectVol: `BACDE`,
		expectChp: `dabcefg`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`ABBBCDE`),
		from:      0,
		to:        1,
		expectVol: `BACDE`,
		expectChp: `bcdaefg`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`BBBBCDE`),
		from:      0,
		to:        1,
		expectVol: `BACDE`,
		expectChp: `abcdefg`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`AAAACDE`),
		from:      0,
		to:        1,
		expectVol: `BACDE`,
		expectChp: `abcdefg`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`ABCDE`),
		from:      1,
		to:        2,
		expectVol: `ACBDE`,
		expectChp: `acbde`,
	}, {
		name:      `move right adjacent`,
		toc:       setupToc(`ABCDE`),
		from:      3,
		to:        4,
		expectVol: `ABCED`,
		expectChp: `abced`,
	}, {
		name:      `move left`,
		toc:       setupToc(`ABCDE`),
		from:      2,
		to:        0,
		expectVol: `CABDE`,
		expectChp: `cabde`,
	}, {
		name:      `move left`,
		toc:       setupToc(`ABCDE`),
		from:      4,
		to:        0,
		expectVol: `EABCD`,
		expectChp: `eabcd`,
	}, {
		name:      `move left adjacent`,
		toc:       setupToc(`ABCDE`),
		from:      1,
		to:        0,
		expectVol: `BACDE`,
		expectChp: `bacde`,
	}, {
		name:      `move left adjacent`,
		toc:       setupToc(`ABCDE`),
		from:      3,
		to:        4,
		expectVol: `ABCED`,
		expectChp: `abced`,
	}}

	testcases := []struct {
		xtest.Base
		toc TOC
	}{{
		Base: xtest.New(`443Q65`),
		toc:  setupToc(``),
	}}
	for _, tc := range testcases {
		t.Run(tc.String(), func(t *testing.T) {})
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf(`%d_%s`, i, tc.name), func(t *testing.T) {
			id, seq := tc.toc.volumes[tc.from].ID, tc.to
			toc, err := tc.toc.ReorderVolume(id, seq)
			if err != nil {
				t.Errorf(`unexpected error: %s`, err)
			}

			if toc.volTitles() != tc.expectVol {
				t.Errorf(`unexpected sequence: %s`, tc.toc.volTitles())
			}

			if toc.chpTitles() != tc.expectChp {
				t.Errorf(`unexpected sequence: %s`, tc.toc.chpTitles())
			}
		})
	}
}

func TestTOCReorderVolumeError(t *testing.T) {
	id1, id2, id3, id4, id5 := ulid.Make(), ulid.Make(), ulid.Make(), ulid.Make(), ulid.Make()
	vol5 := func() []Volume {
		return []Volume{
			{ID: id1, Title: "A"},
			{ID: id2, Title: "B"},
			{ID: id3, Title: "C"},
			{ID: id4, Title: "D"},
			{ID: id5, Title: "E"},
		}
	}
	id6 := ulid.Make()

	tests := []struct {
		name string
		toc  TOC
		id   ID
		seq  int
		err  string
	}{{
		name: `out of range`,
		toc:  TOC{volumes: vol5()},
		id:   id5,
		seq:  -1,
		err:  `out of range`,
	}, {
		name: `out of range`,
		toc:  TOC{volumes: vol5()},
		id:   id5,
		seq:  5,
		err:  `out of range`,
	}, {
		name: `not found`,
		toc:  TOC{volumes: vol5()},
		id:   id6,
		seq:  3,
		err:  `not found`,
	}}

	for i, tc := range tests {
		t.Run(fmt.Sprintf(`%d_%s`, i, tc.name), func(t *testing.T) {
			_, err := tc.toc.ReorderVolume(tc.id, tc.seq)
			if err == nil {
				t.Errorf(`expecting error`)
			}

			if !strings.Contains(err.Error(), tc.err) {
				t.Errorf(`expecting error contains %s, actual: %s`, tc.err, err)
			}
			t.Log(err)
		})
	}
}

func TestTOC_DropVolume(t *testing.T) {
	id1, _, _, _, _, vol := vol5()

	tests := []struct {
		name   string
		toc    TOC
		id     ID
		expect string
	}{{
		name: `basic`,
		toc: TOC{
			volumes: vol,
		},
		id:     id1,
		expect: `BCDE`,
	}}

	for i, tc := range tests {
		t.Run(fmt.Sprintf(`%d_%s`, i, tc.name), func(t *testing.T) {
			toc, err := tc.toc.DropVolume(tc.id)
			if err != nil {
				t.Error(err)
			}

			if toc.volTitles() != tc.expect {
				t.Errorf(`unexpected sequence: %s`, toc.volTitles())
			}
		})
	}
}

func vol5() (id1, id2, id3, id4, id5 ulid.ULID, vol []Volume) {
	id1, id2, id3, id4, id5 = ulid.Make(), ulid.Make(), ulid.Make(), ulid.Make(), ulid.Make()
	vol = []Volume{
		{ID: id1, Title: "A"},
		{ID: id2, Title: "B"},
		{ID: id3, Title: "C"},
		{ID: id4, Title: "D"},
		{ID: id5, Title: "E"},
	}

	return id1, id2, id3, id4, id5, vol
}

func (t TOC) volTitles() string {
	var sb strings.Builder
	for i := range t.volumes {
		sb.WriteString(t.volumes[i].Title)
	}

	return sb.String()
}

func (t TOC) chpTitles() string {
	var sb strings.Builder
	for i := range t.chapters {
		sb.WriteString(t.chapters[i].title)
	}

	return sb.String()
}

func chp5() (id1, id2, id3, id4, id5 ulid.ULID, chp []Chapter) {
	id1, id2, id3, id4, id5 = ulid.Make(),
		ulid.Make(),
		ulid.Make(),
		ulid.Make(),
		ulid.Make()

	chp = []Chapter{
		{ID: id1, title: `a`},
		{ID: id2, title: `b`},
		{ID: id3, title: `c`},
		{ID: id4, title: `d`},
		{ID: id5, title: `e`},
	}

	return id1, id2, id3, id4, id5, chp
}

func TestTOC_ChapterSequenceInVolume(t *testing.T) {
	id1, id2, id3, _, _, chp := chp5()
	vid1, vid2 := ulid.Make(), ulid.Make()
	chp[0].VolumeID = vid1
	chp[1].VolumeID = vid2
	chp[2].VolumeID = vid1
	chp[3].VolumeID = vid2
	chp[4].VolumeID = vid1

	tests := []struct {
		desc   string
		toc    TOC
		id     ID
		expect int
	}{{
		desc:   ``,
		toc:    TOC{chapters: chp},
		id:     id1,
		expect: 1,
	}, {
		desc:   ``,
		toc:    TOC{chapters: chp},
		id:     id3,
		expect: 2,
	}, {
		desc:   ``,
		toc:    TOC{chapters: chp},
		id:     id2,
		expect: 1,
	}, {
		desc:   `not exist`,
		toc:    TOC{chapters: chp},
		id:     ulid.Zero,
		expect: 0,
	}}

	for i, tc := range tests {
		t.Run(fmt.Sprintf(`%d_%s`, i, tc.desc), func(t *testing.T) {
			n := tc.toc.ChapterSequenceInVolume(tc.id)
			if n != tc.expect {
				t.Errorf(`expect: %d, actual: %d`, tc.expect, n)
			}
		})
	}
}

func TestTOC_chaptersSortedByVolumeID(t *testing.T) {
	vid1, vid2, vid3, vid4, vid5, vol := vol5()
	_, _, _, _, _, chp0 := chp5()
	_, _, _, _, _, chp1 := chp5()
	chp := append(chp0, chp1...)

	resetChp := func() {
		for i := range chp {
			chp[i].VolumeID = ulid.Zero
		}
	}

	t.Run(`1`, func(t *testing.T) {
		resetChp()
		lo.ForEach(chp, func(_ Chapter, i int) { chp[i].VolumeID = vid1 })
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err != nil {
			t.Error(err)
		}
	})

	t.Run(`normal_222_33_5`, func(t *testing.T) {
		resetChp()

		chp[0].VolumeID = vid2
		chp[1].VolumeID = vid2
		chp[2].VolumeID = vid2

		chp[3].VolumeID = vid3
		chp[4].VolumeID = vid3

		for i := range chp[5:] {
			chp[5+i].VolumeID = vid5
		}

		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err != nil {
			t.Error(err)
		}
	})

	t.Run(`232`, func(t *testing.T) {
		resetChp()

		chp[0].VolumeID = vid2
		chp[1].VolumeID = vid3
		lo.ForEach(chp[2:], func(_ Chapter, i int) { chp[i+2].VolumeID = vid2 })
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err == nil {
			t.Error(err)
		}
	})

	t.Run(`x4`, func(t *testing.T) {
		resetChp()
		chp[0].VolumeID = ulid.Make()
		lo.ForEach(chp[1:], func(_ Chapter, i int) { chp[i+1].VolumeID = vid4 })
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err == nil {
			t.Error(err)
		}
	})

	t.Run(`4x`, func(t *testing.T) {
		resetChp()
		lo.ForEach(chp, func(_ Chapter, i int) { chp[i].VolumeID = vid4 })
		chp[9].VolumeID = ulid.Make()
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err == nil {
			t.Error(err)
		}
	})

	t.Run(`21`, func(t *testing.T) {
		resetChp()
		chp[0].VolumeID = vid2
		lo.ForEach(chp[1:], func(_ Chapter, i int) { chp[i+1].VolumeID = vid1 })
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err == nil {
			t.Error(err)
		}
	})

	t.Run(`221`, func(t *testing.T) {
		resetChp()
		chp[0].VolumeID = vid2
		chp[1].VolumeID = vid2
		lo.ForEach(chp[2:], func(_ Chapter, i int) { chp[i+2].VolumeID = vid1 })
		toc := TOC{chapters: chp, volumes: vol}
		if err := toc.chaptersSortedByVolumeID(); err == nil {
			t.Error(err)
		}
	})
}

func TestTOC_chapterRangeByVolumeIndex(t *testing.T) {
	testcases := []struct {
		xtest.Base
		toc              TOC
		vi               int
		expectL, expectR int
	}{{
		Base:    xtest.New(`JP9F7R`),
		toc:     setupToc(``),
		vi:      0,
		expectL: 0,
		expectR: 0,
	}, {
		Base:    xtest.New(`CWM3PS`),
		toc:     setupToc(`A`),
		vi:      0,
		expectL: 0,
		expectR: 1,
	}, {
		Base:    xtest.New(`YCMW7Z`),
		toc:     setupToc(`AAA`),
		vi:      0,
		expectL: 0,
		expectR: 3,
	}, {
		Base:    xtest.New(`YCMQXB`),
		toc:     setupToc(`AAA`),
		vi:      -1,
		expectL: -1,
		expectR: -1,
	}, {
		Base:    xtest.New(`YCMWMP`),
		toc:     setupToc(`AAA`),
		vi:      2,
		expectL: -1,
		expectR: -1,
	}, {
		Base:    xtest.New(`8K2HD6`),
		toc:     setupToc(`ABBB`),
		vi:      0,
		expectL: 0,
		expectR: 1,
	}, {
		Base:    xtest.New(`8K2QSZ`),
		toc:     setupToc(`ABBBCC`),
		vi:      1,
		expectL: 1,
		expectR: 4,
	}, {
		Base:    xtest.New(`2YLEE2`),
		toc:     setupToc(`B`),
		vi:      0,
		expectL: 0,
		expectR: 0,
	}, {
		Base:    xtest.New(`2YLFJB`),
		toc:     setupToc(`B`),
		vi:      1,
		expectL: 0,
		expectR: 1,
	}, {
		Base:    xtest.New(`2YLZE2`),
		toc:     setupToc(`B`),
		vi:      2,
		expectL: -1,
		expectR: -1,
	}, {
		Base:    xtest.New(`CB9NP9`),
		toc:     setupToc(`C`),
		vi:      1,
		expectL: 0,
		expectR: 0,
	}, {
		Base:    xtest.New(`9YVH5B`),
		toc:     setupToc(`C`),
		vi:      2,
		expectL: 0,
		expectR: 1,
	}, {
		Base:    xtest.New(`82BUGX`),
		toc:     setupToc(`CDD`),
		vi:      3,
		expectL: 1,
		expectR: 3,
	}}

	for _, tc := range testcases {
		t.Run(tc.String(), func(t *testing.T) {
			l, r := tc.toc.chapterRangeByVolumeIndex(tc.vi)
			if l != tc.expectL || r != tc.expectR {
				t.Errorf(
					`expect: %d %d, actual: %d %d`,
					tc.expectL,
					tc.expectR,
					l,
					r,
				)
			}
		})
	}
}

func TestTOC_ReorderChapter(t *testing.T) {
	testcases := []struct {
		xtest.Base
		toc TOC
		// ci is absolute sequence, seq is target relative sequence.
		ci, seq int
		expect  string
	}{{
		Base:   xtest.New(`4U3G24`, `noop`),
		toc:    setupToc(`A`),
		ci:     0,
		seq:    0,
		expect: `a`,
	}, {
		Base:   xtest.New(`BJC8GS`, `noop exceed`),
		toc:    setupToc(`ABC`),
		ci:     1,
		seq:    1,
		expect: `abc`,
	}, {
		Base:   xtest.New(`M23HZ8`, `noop exceed`),
		toc:    setupToc(`ABC`),
		ci:     2,
		seq:    2,
		expect: `abc`,
	}, {
		Base:   xtest.New(`XM2Y3P`, `left to right exceed`),
		toc:    setupToc(`ABC`),
		ci:     2,
		seq:    3,
		expect: `abc`,
	}, {
		Base:   xtest.New(`7793DP`, `left to right`),
		toc:    setupToc(`AAAABC`),
		ci:     0,
		seq:    2,
		expect: `bcadef`,
	}, {
		Base:   xtest.New(`ZU93DP`, `left to right`),
		toc:    setupToc(`AAAABC`),
		ci:     0,
		seq:    3,
		expect: `bcdaef`,
	}, {
		Base:   xtest.New(`ZUE8DP`, `left to right exceed`),
		toc:    setupToc(`AAAABC`),
		ci:     0,
		seq:    4,
		expect: `bcdaef`,
	}, {
		Base:   xtest.New(`EYX2Y3`, `left to right exceed`),
		toc:    setupToc(`AAAABC`),
		ci:     2,
		seq:    4,
		expect: `abdcef`,
	}, {
		Base:   xtest.New(`EYEK83`, `right to left noop`),
		toc:    setupToc(`BC`),
		ci:     1,
		seq:    2,
		expect: `ab`,
	}, {
		Base:   xtest.New(`MQ2EC2`, `right to left`),
		toc:    setupToc(`BCCCCD`),
		ci:     4,
		seq:    2,
		expect: `abcedf`,
	}}

	for _, tc := range testcases {
		t.Run(tc.String(), func(t *testing.T) {
			cid := tc.toc.chapters[tc.ci].ID
			toc, err := tc.toc.ReorderChapter(cid, tc.seq)
			if err != nil {
				t.Errorf(`unexpected err: %s`, err)
			}

			if toc.chpTitles() != tc.expect {
				t.Errorf(`expected: %s, actual: %s`, tc.expect, toc.chpTitles())
			}
		})
	}
}

func TestTOC_MoveToVolume(t *testing.T) {
	t.Run(`errors`, func(t *testing.T) {
		toc := setupToc(`ABC`)
		fakeID := ulid.Make()

		t.Run(`chapter not found`, func(t *testing.T) {
			_, err := toc.MoveToVolume(fakeID, toc.volumes[0].ID)
			if err == nil || !strings.Contains(err.Error(), `not found`) {
				t.Errorf(`expected chapter not found, got: %v`, err)
			}
		})

		t.Run(`volume not found`, func(t *testing.T) {
			_, err := toc.MoveToVolume(toc.chapters[0].ID, fakeID)
			if err == nil || !strings.Contains(err.Error(), `not found`) {
				t.Errorf(`expected volume not found, got: %v`, err)
			}
		})
	})

	// success cases use chpTitles() to verify final chapter order.
	tests := []struct {
		xtest.Base
		pattern string
		ci      int    // chapter index to move
		target  int    // target volume index
		expect  string // expected chpTitles
	}{
		{
			// Noop — chapter already in target volume
			Base:    xtest.New(`K43P2N`, `noop single vol`),
			pattern: `A`,
			ci:      0,
			target:  0,
			expect:  `a`,
		},
		{
			Base:    xtest.New(`K43Q2G`, `noop already in target first`),
			pattern: `AABBC`,
			ci:      0,
			target:  0,
			expect:  `abcde`,
		},
		{
			Base:    xtest.New(`K43R7H`, `noop already in target last`),
			pattern: `AABBC`,
			ci:      1,
			target:  0,
			expect:  `abcde`,
		},
		{
			Base:    xtest.New(`L1`, `adjacent single-ch vol`),
			pattern: `ABCD`,
			ci:      0,
			target:  1,
			expect:  `bacd`,
		},
		{
			Base:    xtest.New(`L2`, `non-adjacent`),
			pattern: `ABCDE`,
			ci:      0,
			target:  2,
			expect:  `bcade`,
		},
		{
			Base:    xtest.New(`L3`, `multi-ch source first ch`),
			pattern: `AABC`,
			ci:      0,
			target:  1,
			expect:  `bcad`,
		},
		{
			Base:    xtest.New(`L4`, `multi-ch source non-first`),
			pattern: `AABC`,
			ci:      1,
			target:  1,
			expect:  `acbd`,
		},
		{
			Base:    xtest.New(`L5`, `multi-ch target`),
			pattern: `ABB`,
			ci:      0,
			target:  1,
			expect:  `bca`,
		},
		{
			Base:    xtest.New(`L6`, `source becomes empty`),
			pattern: `ABC`,
			ci:      0,
			target:  2,
			expect:  `bca`,
		},
		{
			Base:    xtest.New(`R1`, `adjacent single-ch vol`),
			pattern: `ABCD`,
			ci:      2,
			target:  1,
			expect:  `acbd`,
		},
		{
			Base:    xtest.New(`R2`, `non-adjacent`),
			pattern: `ABCDE`,
			ci:      4,
			target:  1,
			expect:  `aebcd`,
		},
		{
			Base:    xtest.New(`R3`, `multi-ch source move last`),
			pattern: `AABBC`,
			ci:      3,
			target:  0,
			expect:  `adbce`,
		},
		{
			Base:    xtest.New(`R4`, `end to front`),
			pattern: `ABCDE`,
			ci:      4,
			target:  0,
			expect:  `eabcd`,
		},
		{
			Base:    xtest.New(`B1`, `ci == l-1 left to right`),
			pattern: `ABCD`,
			ci:      1,
			target:  2,
			expect:  `acbd`,
		},
		// Empty target volume
		{
			Base:    xtest.New(`E1`, `ci < l to empty vol`),
			pattern: `AAC`,
			ci:      0,
			target:  1,
			expect:  `bac`,
		},
		{
			Base:    xtest.New(`E2`, `ci >= r to empty vol`),
			pattern: `AAC`,
			ci:      2,
			target:  1,
			expect:  `abc`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.String(), func(t *testing.T) {
			toc := setupToc(tc.pattern)
			cid := toc.chapters[tc.ci].ID
			vid := toc.volumes[tc.target].ID

			result, err := toc.MoveToVolume(cid, vid)
			if err != nil {
				t.Fatalf(`unexpected error: %v`, err)
			}

			if title := result.chpTitles(); title != tc.expect {
				t.Errorf(`expected chapters %q, got %q`, tc.expect, title)
			}
		})
	}
}
