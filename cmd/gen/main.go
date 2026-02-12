package main

import (
	"context"
	"fmt"
	"github.com/magomedcoder/gen/api/pb"
	"github.com/magomedcoder/gen/config"
	"github.com/magomedcoder/gen/internal/handler"
	"github.com/magomedcoder/gen/internal/repository"
	"github.com/magomedcoder/gen/internal/repository/postgres"
	"github.com/magomedcoder/gen/internal/service"
	"github.com/magomedcoder/gen/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close(ctx)

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	sessionRepo := postgres.NewChatSessionRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	jwtService := service.NewJWTService(cfg)
	llmRepo, err := repository.NewLLMRunnerRepository(cfg.LLMRunner.Address, cfg.LLMRunner.Model)
	if err != nil {
		log.Fatalf("Ошибка подключения к llm-runner: %v", err)
	}
	defer llmRepo.Close()

	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	chatUseCase := usecase.NewChatUseCase(sessionRepo, messageRepo, llmRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, tokenRepo, jwtService)

	authHandler := handler.NewAuthHandler(authUseCase)
	chatHandler := handler.NewChatHandler(chatUseCase, authUseCase)
	userHandler := handler.NewUserHandler(userUseCase, authUseCase)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	pb.RegisterChatServiceServer(grpcServer, chatHandler)
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера на адресе %s: %v", addr, err)
	}

	log.Printf("запущен на %s", addr)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Ошибка работы сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()
	log.Println("Сервер остановлен")
}
