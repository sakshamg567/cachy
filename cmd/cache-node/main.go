package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/sakshamg567/cachy/internal/cache"
	"github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "port to run cache node on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	grpcServer := grpc.NewServer()
	node := cache.NewCacheNode(100)
	cacheNodepb.RegisterCacheServer(grpcServer, node)
	log.Printf("cache node running on port %s", *port)
	log.Fatal(grpcServer.Serve(lis))
}
