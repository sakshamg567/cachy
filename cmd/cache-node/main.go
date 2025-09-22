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

	// Add a distinctive prefix; keep standard flags (date/time)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix(fmt.Sprintf("[cache-%s] ", *port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	grpcServer := grpc.NewServer()
	node := cache.NewCacheNode(100)
	cacheNodepb.RegisterCacheServer(grpcServer, node)
	log.Printf("starting capacity=%d", 100)
	log.Fatal(grpcServer.Serve(lis))
}
