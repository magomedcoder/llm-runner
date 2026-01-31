package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	runner "github.com/magomedcoder/llm-runner"
	"github.com/magomedcoder/llm-runner/config"
	"github.com/magomedcoder/llm-runner/logger"
	"github.com/magomedcoder/llm-runner/pb"
	"github.com/magomedcoder/llm-runner/provider"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.E("config: %v", err)
		os.Exit(1)
	}
	logger.Default.SetLevel(logger.ParseLevel(cfg.Log.Level))
	textProvider, err := provider.NewTextProvider(cfg)
	if err != nil {
		logger.E("text provider: %v", err)
		os.Exit(1)
	}

	runnerServer := runner.NewServer(textProvider, cfg.MaxConcurrentGenerations)
	lis, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		logger.E("listen: %v", err)
		os.Exit(1)
	}

	defer lis.Close()

	grpcServer := grpc.NewServer()

	pb.RegisterLLMRunnerServiceServer(grpcServer, runnerServer)
	go func() {
		logger.I("listening on %s", cfg.ListenAddr)
		_ = grpcServer.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	grpcServer.GracefulStop()
}
