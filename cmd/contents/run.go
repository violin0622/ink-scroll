package contents

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"play-ddd/contents/app"
	"play-ddd/contents/domain/chapter"
	"play-ddd/contents/domain/novel"
	"play-ddd/contents/iface/grpc"
	ebstderr "play-ddd/contents/infra/eventbus/stderr"
	"play-ddd/contents/infra/outbox"
	"play-ddd/contents/infra/repository/pg"
	"play-ddd/contents/infra/server"
	"play-ddd/contents/infra/server/requestlimiter/bps"
	"play-ddd/contents/infra/server/requestlimiter/composite"
	"play-ddd/contents/infra/server/requestlimiter/nop"
	"play-ddd/contents/infra/server/requestlimiter/qps"
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

	apiName := `/play-ddd/contents`
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
