package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/anzboi/basic-grpc/pkg/svc"
)

type Svc struct{}

// Greet (Unary Client and Server) returns a greeting to the name in the request
func (s Svc) Greet(ctx context.Context, req *svc.GreetRequest) (*svc.GreetResponse, error) {
	return &svc.GreetResponse{Greeting: "Hello " + req.GetName()}, nil
}

// Chat (Streaming client and server) echoes every request message it recieves
func (s Svc) Chat(chatServer svc.Service_ChatServer) error {
	log.Println("Beginning chat stream")
	for {
		m, err := chatServer.Recv()
		if err == io.EOF {
			log.Println("Stream closed by caller")
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("recieved:", m.GetMessage())
		resp := &svc.MessageResponse{
			Response: m.GetMessage(),
		}
		log.Println("responding:", resp.GetResponse())
		if err := chatServer.Send(resp); err != nil {
			log.Println(err)
		}
	}
}

// Echo (Unary Client, Streaming Server) echoes the request every 5 seconds until the client stops listening
func (s Svc) Echo(req *svc.EchoRequest, echoServer svc.Service_EchoServer) error {
	log.Println("Beginning echo")
	echoServer.Send(&svc.EchoResponse{Echo: req.GetMessage()})
	ticker := time.NewTicker(5 * time.Second)

echo:
	for {
		select {
		case <-ticker.C:
			if err := echoServer.Send(&svc.EchoResponse{Echo: req.GetMessage()}); err != nil {
				log.Println(err)
				break echo
			}
		case <-echoServer.Context().Done():
			log.Println("client cancelled the request")
			break echo
		}
	}

	ticker.Stop()
	return nil
}
