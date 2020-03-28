package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/anzboi/basic-grpc/pkg/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("no PORT env variable, defaulting to 8080")
		port = "8080"
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	svr := grpc.NewServer()
	svc.RegisterServiceServer(svr, Svc{})
	reflection.Register(svr)
	log.Printf("listening on port %s", port)
	log.Fatal(svr.Serve(lis))
}
