package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	client "lyceum_service/pkg/api/order"
	"lyceum_service/pkg/logger"
	"net"
	"net/http"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
}

func New(ctx context.Context, port, restPort int, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	client.RegisterOrderServiceServer(grpcServer, NewOrderService(service))

	restSrv := runtime.NewServeMux()
	if err := client.RegisterOrderServiceHandlerServer(context.Background(), restSrv, NewOrderService(service)); err != nil {
		return nil, err
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: restSrv,
	}
	return &Server{grpcServer, httpServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting gRPC server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
		return s.grpcServer.Serve(s.listener)
	})

	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting Rest server", zap.String("port", s.restServer.Addr))
		return s.restServer.ListenAndServe()
	})

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "gRPC server stopped")
	}
	return s.restServer.Shutdown(ctx)
}
