package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	runner "github.com/magomedcoder/llm-runner"
	"github.com/magomedcoder/llm-runner/config"
	"github.com/magomedcoder/llm-runner/gpu"
	"github.com/magomedcoder/llm-runner/logger"
	"github.com/magomedcoder/llm-runner/pb/llmrunnerpb"
	"github.com/magomedcoder/llm-runner/provider"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func cmdServe() *cli.Command {
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Запустить gRPC-сервер раннера (конфиг из LLM_RUNNER_CONFIG / по умолчанию)",
		Action:  runServe,
	}
}

func runServe(ctx context.Context, _ *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		logger.Default.SetLevel(logger.LevelInfo)
		logger.E("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	logger.Default.SetLevel(logger.ParseLevel(cfg.LogLevel))
	logger.I("Запуск раннера")

	textProvider, err := provider.NewTextProvider(cfg)
	if err != nil {
		logger.E("Движок текста: %v", err)
		return err
	}

	warmCtx, warmCancel := context.WithTimeout(ctx, 30*time.Minute)
	defer warmCancel()
	if err := textProvider.WarmDefaultModel(warmCtx, cfg.DefaultModel); err != nil {
		logger.E("Загрузка default_model %q: %v", cfg.DefaultModel, err)
		return fmt.Errorf("не удалось загрузить default_model: %w", err)
	}
	logger.I("Модель по умолчанию загружена: %s", strings.TrimSpace(cfg.DefaultModel))

	gpuCollector := gpu.NewCollector()
	runnerServer := runner.NewServer(textProvider, gpuCollector, cfg.MaxConcurrentGenerations, cfg.DefaultModel, cfg.UnloadAfterRPC(), cfg.ModelsDir())

	listenAddr := cfg.ListenAddr()
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.E("Ошибка слушателя: %v", err)
		return err
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(1024),
	)
	llmrunnerpb.RegisterLLMRunnerServiceServer(grpcServer, runnerServer)

	go func() {
		logger.I("Раннер слушает на %s", listenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			logger.E("Ошибка gRPC: %v", err)
			os.Exit(1)
		}
	}()

	coreAddr := cfg.CoreAddr()
	if coreAddr != "" && listenAddr != "" {
		if err := registerWithCore(coreAddr, listenAddr, cfg.RegistrationToken); err != nil {
			logger.W("Регистрация в ядре не удалась: %v", err)
		} else {
			logger.I("Зарегистрирован в ядре %s как %s", coreAddr, listenAddr)
			defer unregisterFromCore(coreAddr, listenAddr, cfg.RegistrationToken)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-ctx.Done():
	}

	grpcServer.GracefulStop()
	logger.I("Раннер остановлен")

	return nil
}

func registerWithCore(coreAddr, registerAddress, registrationToken string) error {
	conn, err := grpc.NewClient(coreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("подключение к ядру: %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := llmrunnerpb.NewLLMRunnerServiceClient(conn)
	_, err = client.RegisterRunnerWithToken(ctx, &llmrunnerpb.RunnerRegisterRequest{
		ListenAddress:     registerAddress,
		RegistrationToken: registrationToken,
	})

	return err
}

func unregisterFromCore(coreAddr, registerAddress, registrationToken string) {
	conn, err := grpc.NewClient(coreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.W("Unregister: подключение к ядру: %v", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := llmrunnerpb.NewLLMRunnerServiceClient(conn)
	_, _ = client.UnregisterRunnerWithToken(ctx, &llmrunnerpb.RunnerUnregisterRequest{
		ListenAddress:     registerAddress,
		RegistrationToken: registrationToken,
	})
}
