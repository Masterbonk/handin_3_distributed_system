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

func addToLamport(inputLamport int32, ourLamport *int32) {
	if inputLamport > *ourLamport {
		*ourLamport = inputLamport
	}
	*ourLamport++
}

func main() {
	var username string

	var lamport int32 = 0

	flag.StringVar(&username, "u", "anonymous", "Set the client username. Defaults to anonymous")
	flag.Parse()

	testMessages := []*cc.ClientMessage{
		{Msg: "First message", ClientName: username},
		{Msg: "Second message", ClientName: username},
		{Msg: "Third message", ClientName: username},
		{Msg: "Fourth message", ClientName: username},
		{Msg: "Fifth message", ClientName: username},
		{Msg: "the message that is 128 charachters long. the message that is 128 charachters long. the message that is 128 charachters long. the message that is 128 charachters long. ", ClientName: username},

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
	go func(lamport *int32) {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			} else{
				addToLamport(in.Lamport, lamport)
			}

			if in.MessageType == 0 {
				fmt.Printf("Time: %d, %s\n",*lamport, in.Msg)
			} else if in.MessageType == 1 {
				fmt.Printf("Time: %d, %s: %s\n",*lamport, in.ClientName, in.Msg)
			}
		}
	}(&lamport)

	for _, msg := range testMessages {
		err := stream.Send(msg)
		lamport++
		if err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", msg, err)
		}
		time.Sleep(time.Second * 3)
	}
	stream.CloseSend()
	fmt.Println("Close")
	<-waitc
}
