package app

import (
	"fmt"
	"net"
	"sync"

	"github.com/lgu-elo/auth/internal/adapters/database"
	"github.com/lgu-elo/auth/internal/auth"
	"github.com/lgu-elo/auth/internal/config"
	"github.com/lgu-elo/auth/internal/server"
	"github.com/lgu-elo/auth/pkg/pb"
	"github.com/lgu-elo/gateway/pkg/logger"
	"github.com/lgu-elo/gateway/pkg/rpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	fxlogrus "github.com/takt-corp/fx-logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.WithLogger(func(log *logrus.Logger) fxevent.Logger {
			return &fxlogrus.LogrusLogger{Logger: log}
		}),
		fx.Provide(
			logger.New,
			config.New,
			database.New,
			func() *sync.Mutex {
				var mu sync.Mutex
				return &mu
			},

			fx.Annotate(auth.NewStorage, fx.As(new(auth.Repository))),
			fx.Annotate(auth.NewService, fx.As(new(auth.IService))),
			fx.Annotate(auth.NewHandler, fx.As(new(auth.IHandler))),

			server.NewAPI,

			func(logger *logrus.Logger) *grpc.Server {
				return rpc.New(
					rpc.WithLoggerLogrus(logger),
				)
			},
		),
		fx.Invoke(serve),
	)
}

func serve(logger *logrus.Logger, srv *grpc.Server, cfg *config.Cfg, api *server.API) {
	pb.RegisterAuthServer(srv, api)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		logger.Error(errors.Wrap(err, "cannot bind server"))
		return // или панику выбросить
	}

	if errServe := srv.Serve(lis); errServe != nil {
		logger.Error(errors.Wrap(err, "cannot lsiten server"))
		return // или панику выбросить
	}
}
