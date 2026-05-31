package vo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
)

func TestTOCReorderVolume(t *testing.T) {
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

	vol1 := func() []Volume { return []Volume{{ID: id1, Title: "A"}} }

	tests := []struct {
		name   string
		toc    TOC
		id     ID
		seq    int
		expect string
	}{{
		name:   `noop`,
		toc:    TOC{volumes: vol5()},
		id:     id1,
		seq:    0,
		expect: `ABCDE`,
	}, {
		name:   `noop`,
		toc:    TOC{volumes: vol5()},
		id:     id3,
		seq:    2,
		expect: `ABCDE`,
	}, {
		name:   `noop`,
		toc:    TOC{volumes: vol5()},
		id:     id5,
		seq:    4,
		expect: `ABCDE`,
	}, {
		name:   `move backward adjacent`,
		toc:    TOC{volumes: vol5()},
		id:     id2,
		seq:    0,
		expect: `BACDE`,
	}, {
		name:   `move backward`,
		toc:    TOC{volumes: vol5()},
		id:     id1,
		seq:    1,
		expect: `BACDE`,
	}, {
		name:   `move backward adjacent`,
		toc:    TOC{volumes: vol5()},
		id:     id3,
		seq:    1,
		expect: `ACBDE`,
	}, {
		name:   `move backward adjacent`,
		toc:    TOC{volumes: vol5()},
		id:     id5,
		seq:    3,
		expect: `ABCED`,
	}, {
		name:   `move backward`,
		toc:    TOC{volumes: vol5()},
		id:     id3,
		seq:    0,
		expect: `CABDE`,
	}, {
		name:   `move backward`,
		toc:    TOC{volumes: vol5()},
		id:     id5,
		seq:    0,
		expect: `EABCD`,
	}, {
		name:   `move forward adjacent`,
		toc:    TOC{volumes: vol5()},
		id:     id2,
		seq:    0,
		expect: `BACDE`,
	}, {
		name:   `move forward adjacent`,
		toc:    TOC{volumes: vol5()},
		id:     id4,
		seq:    4,
		expect: `ABCED`,
	}, {
		name:   `noop`,
		toc:    TOC{volumes: vol1()},
		id:     id1,
		seq:    0,
		expect: `A`,
	}}

	for i, tc := range tests {
		t.Run(fmt.Sprintf(`%d_%s`, i, tc.name), func(t *testing.T) {
			toc, err := tc.toc.ReorderVolume(tc.id, tc.seq)
			if err != nil {
				t.Errorf(`unexpected error: %s`, err)
			}

			if toc.volTitles() != tc.expect {
				t.Errorf(`unexpected sequence: %s`, tc.toc.volTitles())
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

	return
}

func (t TOC) volTitles() string {
	var sb strings.Builder
	for i := range t.volumes {
		sb.WriteString(t.volumes[i].Title)
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
		{ID: id1},
		{ID: id2},
		{ID: id3},
		{ID: id4},
		{ID: id5},
	}

	return
}

func TestChapterSequenceInVolume(t *testing.T) {
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

func helper(pattern string, chp []Chapter, vids ...ID) {
	pattern = strings.ReplaceAll(pattern, `_`, ``)
	for _, n := range pattern {
		switch n {
		case 'x':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		}
	}

}
