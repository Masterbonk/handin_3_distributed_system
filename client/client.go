package main

import (
	cc "ChittyChat/ChittyChat"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var username string

	flag.StringVar(&username, "u", "anonymous", "Set the client username. Defaults to anonymous")
	flag.Parse()

	testMessages := []*cc.ClientMessage{
		{Msg: "First message", ClientName: username},
		{Msg: "Second message", ClientName: username},
		{Msg: "Third message", ClientName: username},
		{Msg: "Fourth message", ClientName: username},
		{Msg: "Fifth message", ClientName: username},
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	var serverAddress = "localhost:50051"
	conn, err := grpc.NewClient(serverAddress, opts...)

	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	// close connection when function terminates
	defer conn.Close()

	// create client
	client := cc.NewChittyChatClient(conn)

	// create contect, cancel when function terminates
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	clientDeadline := time.Now().Add(time.Duration(20) * time.Second)
	ctx, cancel = context.WithDeadline(ctx, clientDeadline)

	defer cancel()

	stream, err := client.Chat(ctx)
	stream.Send(&cc.ClientMessage{ClientName: username})

	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}

			if in.MessageType == 0 {
				fmt.Printf("%s\n", in.Msg)
			} else if in.MessageType == 1 {
				fmt.Printf("%s: %s\n", in.ClientName, in.Msg)
			}
		}
	}()

	for _, msg := range testMessages {
		err := stream.Send(msg)
		if err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", msg, err)
		}
		time.Sleep(time.Second * 3)
	}
	stream.CloseSend()
	fmt.Println("Close")
	<-waitc
}
