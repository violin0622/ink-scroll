package novel

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg/convert"
	dtulid "github.com/violin0622/ink-scroll/contents/infra/repository/pg/datatypes/ulid"
)

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
	novel := convert.FromDomain[novel.Novel, Novel](n)
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
