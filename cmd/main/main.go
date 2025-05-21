package main

import (
	"context"
	"github.com/Sanchir01/auth-micro/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	defer cancel()

	go func() { env.GRPCSrv.MustStart() }()
	<-ctx.Done()

	env.GRPCSrv.Stop()
}
