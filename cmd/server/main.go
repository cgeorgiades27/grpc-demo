package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/cgeorgiades27/grpc-demo/pkg/xref"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

var (
	port      = flag.Int("port", 50051, "The server port")
	redisAddr = flag.String("redis", "localhost:6379", "redis server address")
	dataPath  = flag.String("data", "./data/random", "init data path")
)

func main() {

	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s := grpc.NewServer()
	server := xref.NewXrefService(context.Background(), rdb)
	server.InitData(*dataPath)
	xref.RegisterXrefServiceServer(s, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
