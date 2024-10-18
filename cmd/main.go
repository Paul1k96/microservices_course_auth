package main

import (
	"log/slog"
	"net"

	"github.com/Paul1k96/microservices_course_auth/internal/config/env"
	"github.com/Paul1k96/microservices_course_auth/internal/user"
	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.Default()

	grpcConfig := env.NewGRPCConfig()

	listen, err := net.Listen("tcp", grpcConfig.GetAddress())
	if err != nil {
		logger.Error("failed to listen", slog.String("error", err.Error()))
		return
	}

	pgConfig := env.NewPGConfig()

	db, err := sqlx.Connect("postgres", pgConfig.GetDSN())
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		return
	}

	userDB := user.NewUserRepository(db)

	userAPIv1 := user.NewUserAPI(logger, userDB)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	user_v1.RegisterUserServer(grpcServer, userAPIv1)

	logger.Info("server listening at", slog.Any("addr", listen.Addr()))

	if err = grpcServer.Serve(listen); err != nil {
		logger.Error("failed to serve", slog.String("error", err.Error()))
		return
	}
}
