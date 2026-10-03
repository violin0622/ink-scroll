package contents

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/iface/grpc"
	ebstderr "github.com/violin0622/ink-scroll/contents/infra/eventbus/stderr"
	"github.com/violin0622/ink-scroll/contents/infra/outbox"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg"
	"github.com/violin0622/ink-scroll/contents/infra/server"
	"github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter/bps"
	"github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter/composite"
	"github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter/nop"
	"github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter/qps"
)

func Run() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `start contents service`,
		RunE:  runErr,
	}
}

func runErr(*cobra.Command, []string) error {
	db, err := pg.InitDB(pg.DSN)
	if err != nil {
		return err
	}

	novelFac := novel.NewFactory(logr.Discard())
	chapterFac := chapter.NewFactory(logr.Discard())
	repo := pg.New(db, novelFac, chapterFac)
	outbox := outbox.NewRelay(
		ebstderr.New(),
		repo.Outbox(),
		outbox.WithLogr(logr.Discard()),
	)

	ch := app.NewCommandHandler(repo, logr.Discard())
	cmdSvc := grpc.NewCmdService(ch)
	qh := app.NewQueryHandler()
	querySvc := grpc.NewQueryService(qh)

	apiName := `/github.com/violin0622/ink-scroll/contents`
	nopLimiter := nop.New()
	qpsLimiter := qps.New(apiName)
	qpsLimiter.Set(100, 100)
	bpsLimiter := bps.New(apiName)
	bpsLimiter.Set(100, 100)

	rl := composite.New(
		&nopLimiter,
		&qpsLimiter,
		&bpsLimiter,
	)

	svr := server.New(
		&cmdSvc,
		&querySvc,
		server.WithRequestLimiter(rl),
	)

	l, err := net.Listen(`tcp`, `:8080`)
	if err != nil {
		return err
	}

	go func() {
		_c := make(chan os.Signal, 1)
		signal.Notify(_c,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		msg := <-_c
		fmt.Println(`Quit by signal: `, msg.String())
		svr.GracefulStop()
	}()

	if err := outbox.Start(); err != nil {
		fmt.Println(`Failed to start relay.`)
	}

	if err := svr.Serve(l); err != nil {
		fmt.Println(`Server stopped.`)
	}

	_ = outbox.Stop()

	return nil
}
