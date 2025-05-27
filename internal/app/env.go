package app

import (
	grpcapp "github.com/Sanchir01/auth-micro/internal/app/grpc"
	"github.com/Sanchir01/auth-micro/internal/config"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
	DB      *Database
	Log     *slog.Logger
	Config  *config.Config
}

func NewEnv() (*App, error) {
	cfg := config.InitConfig()

	lg := SetupLogger(cfg.Env)
	db, err := NewDataBases(cfg)
	if err != nil {
		lg.Error("pgx error connect", err.Error())
		return nil, err
	}
	repos := NewRepository(db)
	serv := NewServices(repos, db)
	gRPCServer := grpcapp.New(lg, ":"+cfg.GRPC.Port, serv.UserService)
	return &App{
		GRPCSrv: gRPCServer,
		DB:      db,
		Log:     lg,
		Config:  cfg,
	}, nil
}
