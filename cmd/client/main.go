package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("addr", "localhost:50055", "server address")
)

func main() {

	flag.Parse()
	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		log.Fatalf("failed to dial host: %v", err)
	}
	defer conn.Close()
}
