package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anzboi/basic-grpc/pkg/svc"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := svc.NewServiceClient(cc)
	resp, err := client.Greet(context.Background(), &svc.GreetRequest{Name: "Dan"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.GetGreeting())

	ctx, cancel := context.WithCancel(context.Background())
	chatClient, err := client.Chat(ctx)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("send: How do you do")
			if err := chatClient.Send(&svc.Message{Message: "How do you do"}); err != nil {
				log.Printf("SEND error: %v", err)
				break
			}
		}
	}()

	go func() {
		for {
			resp, err := chatClient.Recv()
			if err != nil {
				log.Println("RECV error: %v", err)
				break
			}
			fmt.Println("respose:", resp.GetResponse())
		}
	}()

	time.Sleep(5 * time.Second)
	chatClient.CloseSend()
	cancel()
}
