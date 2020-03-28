package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/anzboi/basic-grpc/pkg/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()
	port := *port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	svr := grpc.NewServer()
	svc.RegisterServiceServer(svr, Svc{})
	reflection.Register(svr)
	log.Printf("listening on port %d", port)
	log.Fatal(svr.Serve(lis))
}
