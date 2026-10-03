package pg_test

import (
	"testing"

	"github.com/violin0622/ink-scroll/contents/infra/repository/pg"
)

func TestInitDB(t *testing.T) {
	_, err := pg.InitDB(pg.DSN)
	if err != nil {
		t.Error(err)
	}
}
