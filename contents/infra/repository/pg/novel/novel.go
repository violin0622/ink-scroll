package novel

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg/convert"
	dt "github.com/violin0622/ink-scroll/contents/infra/repository/pg/datatypes"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg/datatypes/ts"
	dtulid "github.com/violin0622/ink-scroll/contents/infra/repository/pg/datatypes/ulid"
)

// var _ app.NovelRepo = novelRepo{}

type ID = dtulid.ULID

type Novel struct {
	dt.Model[ID]
	AuthorID    ID
	Title       string
	Category    string
	Description string
	Tags        []string  `gorm:"serializer:json"`
	TOC         []Chapter `gorm:"serializer:json"`
	Volumes     []Volume  `gorm:"serializer:json"`
	Status      int
	WordCount   int
}

func (n *Novel) FromDomain(novel novel.Novel) {
	tags := make([]string, len(novel.Tags()))
	for i, t := range novel.Tags() {
		tags[i] = string(t)
	}

	n.ID = dtulid.ULID(novel.ID())
	n.CreatedTs = ts.From(novel.CreatedAt())
	n.UpdatedTs = ts.From(novel.UpdatedAt())
	n.AuthorID = dtulid.ULID(novel.AuthorID())
	n.Title = string(novel.Title())
	n.Description = string(novel.Description())
	n.Category = string(novel.Category())
	n.Status = int(novel.Status())
	n.Tags = tags
	n.TOC = convert.SliceFromDomain[vo.Chapter, Chapter](novel.Chapters())
	n.Volumes = convert.SliceFromDomain[vo.Volume, Volume](novel.Volumes())
	n.WordCount = novel.WordCount()
}

type Chapter struct {
	ID, VolumeID string
}

func (c *Chapter) FromDomain(vc vo.Chapter) {
	if c == nil {
		return
	}

	c.ID = vc.ID.String()
}

func (c *Chapter) IntoDomain() (vc vo.Chapter, err error) {
	if c == nil {
		return vc, nil
	}

	if vc.ID, err = ulid.Parse(c.ID); err != nil {
		return vc, err
	}

	if vc.VolumeID, err = ulid.Parse(c.VolumeID); err != nil {
		return vc, err
	}

	return vc, nil
}

type Volume struct {
	ID    string
	Title string
	Desc  string
	Kind  int
}

func (v *Volume) FromDomain(vv vo.Volume) {
	if v == nil {
		return
	}

	v.Kind = int(vv.Kind)
	v.Title = vv.Title
	v.Desc = vv.Description
	v.ID = vv.ID.String()
}

func (v *Volume) IntoDomain() (vv vo.Volume, err error) {
	if v == nil {
		return vv, err
	}

	if vv.ID, err = ulid.Parse(v.ID); err != nil {
		return vv, err
	}

	vv.Title = v.Title
	vv.Kind = vo.VolumeKind(v.Kind)
	vv.Description = v.Desc

	return vv, err
}
