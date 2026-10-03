package pg

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/infra/outbox"
	pgevent "github.com/violin0622/ink-scroll/contents/infra/repository/pg/event"
	pgnovel "github.com/violin0622/ink-scroll/contents/infra/repository/pg/novel"
)

const DSN = `host=127.0.0.1 port=15432 sslmode=disable user=postgres password=mysecretpassword dbname=playddd`

func InitDB(dsn string) (db *gorm.DB, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		})
	db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		return db, err
	}

	return db, nil
}

var _ app.Repo = postgresRepo{}

func New(
	db *gorm.DB,
	novelFact novel.Factory,
	chapterFac chapter.Factory,
) postgresRepo {
	return postgresRepo{
		db:          db,
		novelFac:    novelFact,
		chapterFact: chapterFac,
	}
}

func (g postgresRepo) fork(db *gorm.DB) postgresRepo {
	return postgresRepo{
		db:       db,
		novelFac: g.novelFac,
	}
}

type postgresRepo struct {
	novelFac    novel.Factory
	chapterFact chapter.Factory

	db *gorm.DB
}

func (g postgresRepo) Event() novel.EventRepo   { return pgevent.New(g.db) }
func (g postgresRepo) Chapter() app.ChapterRepo { panic(`unimplemented`) }
func (g postgresRepo) Outbox() outbox.EventRepo { return pgevent.New(g.db) }

func (g postgresRepo) Novel() app.NovelRepo { return pgnovel.New(g.db, g.novelFac) }
func (g postgresRepo) Tx(fn func(app.Repo) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(g.fork(tx))
	})
}

func G[N any](g postgresRepo, fn func(app.Repo) (N, error)) (N, error) {
	var n N
	return n, g.Tx(func(tx app.Repo) error {
		var err error
		n, err = fn(tx)
		return err
	})
}
