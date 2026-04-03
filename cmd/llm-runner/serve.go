package main

import (
	"context"
	"fmt"
	"net"
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
		Usage:   "Запустить сервер раннера",
		Action:  runServe,
	}
}

func runServe(ctx context.Context, _ *cli.Command) error {
	ctx, stopSignals := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	cfg, err := config.Load()
	if err != nil {
		logger.Default.SetLevel(logger.LevelInfo)
		logger.E("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	logger.Default.SetLevel(logger.ParseLevel(cfg.LogLevel))
	logger.I("Запуск раннера")
	logger.I("Конфиг: gpu_layers=%d (0 - только CPU; -1 - все слои на GPU; нужны make build-gpu и CUDA-сборка libllama)", cfg.GpuLayers)

	textProvider, err := provider.NewTextProvider(cfg)
	if err != nil {
		logger.E("Движок текста: %v", err)
		return err
	}

	warmCtx, warmCancel := context.WithTimeout(ctx, 30*time.Minute)
	defer warmCancel()
	warmModel := strings.TrimSpace(cfg.DefaultModel)

	if err := textProvider.WarmDefaultModel(warmCtx, warmModel); err != nil {
		logger.E("Предзагрузка модели %q: %v", warmModel, err)
		return fmt.Errorf("не удалось загрузить модель при старте: %w", err)
	}

	logger.I("Модель по умолчанию загружена: %s", warmModel)

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

	errCh := make(chan error, 1)
	go func() {
		logger.I("Раннер слушает на %s", listenAddr)
		errCh <- grpcServer.Serve(lis)
	}()

	coreAddr := cfg.CoreAddr()
	registered := false
	if coreAddr != "" && listenAddr != "" {
		if err := registerWithCore(coreAddr, listenAddr, cfg.RegistrationToken, cfg); err != nil {
			logger.W("Регистрация в ядре не удалась: %v", err)
		} else {
			logger.I("Зарегистрирован в ядре %s как %s", coreAddr, listenAddr)
			registered = true
		}
	}

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err == nil || err == grpc.ErrServerStopped {
			return nil
		}
		logger.E("Ошибка gRPC: %v", err)
		return err
	}

	logger.I("Остановка раннера...")

	if registered {
		unregisterFromCore(coreAddr, listenAddr, cfg.RegistrationToken)
	}

	const shutdownTimeout = 15 * time.Second
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		logger.W("GracefulStop не завершился за %s, принудительная остановка", shutdownTimeout)
		grpcServer.Stop()
	}

	unloadCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := textProvider.UnloadModel(unloadCtx); err != nil {
		logger.W("Выгрузка модели при остановке: %v", err)
	}

	logger.I("Раннер остановлен")
	return nil
}

func registerHintsProto(cfg *config.Config) *llmrunnerpb.RunnerRegisterHints {
	return &llmrunnerpb.RunnerRegisterHints{
		MaxContextTokens:                     int32(cfg.MaxContextTokens),
		LlmHistoryMaxMessages:                int32(cfg.LLMHistoryMaxMessages),
		LlmHistorySummarizeDropped:           cfg.LLMHistorySummarizeDropped,
		LlmHistorySummaryMaxInputRunes:       int32(cfg.LLMHistorySummaryMaxInputRunes),
		LlmHistorySummaryModel:               strings.TrimSpace(cfg.LLMHistorySummaryModel),
		LlmHistorySummaryRunnerListenAddress: strings.TrimSpace(cfg.LLMHistorySummaryRunnerListen),
		LlmHistorySummaryCacheEntries:        int32(cfg.LLMHistorySummaryCacheEntries),
		MaxToolInvocationRounds:              int32(cfg.MaxToolInvocationRounds),
	}
}

func registerWithCore(coreAddr, registerAddress, registrationToken string, cfg *config.Config) error {
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
		Hints:             registerHintsProto(cfg),
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
