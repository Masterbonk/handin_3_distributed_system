package main

import (
	proto "Server/grpc"
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
Converting from the proto.proto file has the following rules:
- Everything will be pascal case, so server_thing becomes ServerThing
- Fields will always have uppercase first letter, so t becomes T
*/

type Chitty_databaseServer struct {
	proto.UnimplementedChittyChatServer
	time         int32
	clientId     []int
	nextClientId int
}

func main() {
	server := &Chitty_databaseServer{time: 0, clientId: []int{}, nextClientId: 1}

	server.start_server()
}

func (s *Chitty_databaseServer) start_server() {
	grpcServer := grpc.NewServer()
	listiner, err := net.Listen("top", ":5050")
	if err != nil {
		log.Fatalf("Did not work 1")
	}

	proto.RegisterChittyChatServer(grpcServer, s)

	err = grpcServer.Serve(listiner)

	if err != nil {
		log.Fatalf("Did not work 2")
	}
}

func (s *Chitty_databaseServer) ClientJoin(ctx context.Context, in *proto.ClientHasJoined) (*proto.Empty, error) {
	s.time++

	for _, clientId := range s.clientId {
		EstablishConnection(5050+clientId, 1, clientId, s.time)
	}

	s.clientId = append(s.clientId, s.nextClientId)

	s.nextClientId++

	return (s.nextClientId - 1), nil

}

func (s *Chitty_databaseServer) ClientLeft(ctx context.Context, in *proto.ClientHasLeft) (*proto.Empty, error) {
	s.time++

	for _, clientId := range s.clientId {
		EstablishConnection(5050+clientId, 2, clientId, s.time)
	}

	return &proto.Empty{}, nil

}

func EstablishConnection(port, state, clientId int, time int32) {
	conn, err := grpc.NewClient("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working 3")
	}

	client := proto.NewClientClient(conn)

	if state == 1 { //ClientJoined has been called

		ClientJoinProces(client, clientId, time)

	} else if state == 2 { //ClientLeft has been called
		ClientLeftProces(client, clientId, time)

	}

}

func ClientJoinProces(client proto.ClientClient, clientId int, time int32) {

	_, err := client.ClientJoin(context.Background(),
		&proto.ServerClientHasJoined{T: time,
			Message: "Participant " + strconv.Itoa(clientId) + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(time))})

	if err != nil {
		log.Fatalf("Not working 4")
	}
}

func ClientLeftProces(client proto.ClientClient, clientId int, time int32) {

	_, err := client.ClientLeft(context.Background(),
		&proto.ServerClientHasLeft{T: time,
			Message: "Participant " + strconv.Itoa(clientId) + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(time))})

	if err != nil {
		log.Fatalf("Not working 4")
	}
}
