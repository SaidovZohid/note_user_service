package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SaidovZohid/note_user_service/config"
	pb "github.com/SaidovZohid/note_user_service/genproto/user_service"
	"github.com/SaidovZohid/note_user_service/service"
	"github.com/SaidovZohid/note_user_service/storage"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.New(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	strg := storage.NewStorage(psqlConn)
	inMemory := storage.NewRedisStorage(rdb)

	userService := service.NewUserService(strg, inMemory)

	listen, err := net.Listen("tcp", cfg.GrpcPort)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	reflection.Register(s)
	
	log.Printf("gRPC port started in: %v", cfg.GrpcPort)
	
	if s.Serve(listen); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
