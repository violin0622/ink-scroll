package novel

import (
	"context"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"play-ddd/contents/domain/novel"
	"play-ddd/contents/domain/novel/vo"
	dt "play-ddd/contents/infra/repository/pg/datatypes"
	"play-ddd/contents/infra/repository/pg/datatypes/ts"
	dtulid "play-ddd/contents/infra/repository/pg/datatypes/ulid"
)

// var _ app.NovelRepo = novelRepo{}

type ID = dtulid.ULID

type novelRepo struct {
	db   *gorm.DB
	fact novel.Factory
}

// Get implements app.NovelRepo.
func (p novelRepo) Get(
	ctx context.Context,
	id novel.ID,
) (novel novel.Novel, err error) {
	var n Novel
	err = p.db.
		WithContext(ctx).
		Where(nonDeleted).
		Take(&n, dtulid.ULID(id)).
		Error
	if err != nil {
		return novel, err
	}

	return p.restore(n)
}

// Save implements app.NovelRepo.
func (nr novelRepo) Save(ctx context.Context, n novel.Novel) error {
	novel := fromDomain[novel.Novel, Novel](n)
	return nr.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: `id`}},
			UpdateAll: true,
			Where: clause.Where{
				Exprs: []clause.Expression{
					gorm.Expr(`"novels"."updated_ts"<="excluded"."updated_ts"`),
					nonDeleted,
				},
			},
		}).
		Create(&novel).
		Error
}

var nonDeleted = gorm.Expr(`"novels"."deleted_ts" = 0 `)

func New(
	db *gorm.DB,
	fact novel.Factory,
) novelRepo {
	return novelRepo{db: db, fact: fact}
}

type Novel struct {
	dt.Model[ID]
	AuthorID    ID
	Title       string
	Category    string
	Description string
	Tags        []string  `gorm:"serializer:json"`
	TOC         []Chapter `gorm:"serializer:json"`
	Status      int
	WordCount   int
}

func (n *Novel) fromDomain(novel novel.Novel) {
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
	n.TOC = fromSlice[vo.Chapter, Chapter](novel.Chapters())
	n.WordCount = novel.WordCount()
}

type Chapter struct {
	ID       string
	Sequence int
}

func (c *Chapter) fromDomain(vc vo.Chapter) {
	if c == nil {
		return
	}

	c.ID = vc.ID.String()
	c.Sequence = vc.Sequence
}

func (c *Chapter) IntoDomain() (vc vo.Chapter, err error) {
	if c == nil {
		return vc, nil
	}

	if c.ID != "" {
		vc.ID, err = ulid.Parse(c.ID)
		if err != nil {
			return vc, err
		}
	}
	vc.Sequence = c.Sequence
	return vc, nil
}

func fromDomain[A, M any, P fromPtr[A, M]](a A) (b M) {
	P(&b).fromDomain(a)
	return b
}

// fromPtr 约束类型必须是指向 B 的指针，且该指针类型实现了 fromDomain(A) 方法
type fromPtr[A any, M any] interface {
	fromDomain(A)
	*M
}

func fromSlice[A, M any, Mp fromPtr[A, M]](as []A) (bs []M) {
	if len(as) == 0 {
		return make([]M, 0)
	}
	bs = make([]M, len(as))
	for i := range as {
		Mp(&bs[i]).fromDomain(as[i])
	}
	return bs
}
