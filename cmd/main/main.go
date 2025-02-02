package main

import (
	"context"
	"fmt"
	"lyceum_service/internal/config"
	"lyceum_service/internal/repository"
	service "lyceum_service/internal/service"
	"lyceum_service/internal/transport/grpc"
	"lyceum_service/pkg/db/cache"
	"lyceum_service/pkg/db/postgres"
	"lyceum_service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "lyceum"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg := config.New()
	if cfg == nil {
		panic("failed to load config")
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	redis := cache.New(cfg.RedisConfig)
	fmt.Println(redis.Ping(ctx))

	repo := repository.NewOrderRepository(db)

	srv := service.NewOrderService(repo)

	grpcserver, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort, srv)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcserver.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	if err := grpcserver.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}
