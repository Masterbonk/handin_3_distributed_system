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

type Chitty_databaseServer struct {
	proto.UnimplementedChittyChatServer
	time int32
	clientId []int
}

func main() {
	server := &Chitty_databaseServer{time: 0, clientId: []int{}}

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
	time++
}

func EstablishConnection (port, state, clientId int) {
	conn, err := grpc.NewClient("localhost:" + strconv.Itoa(port) , grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Not working 3")
    }

    client := proto.NewClientClient(conn)

	if state == 1 {

		ClientJoinProces(client, clientId)
		
	}

    
   
}

func ClientJoinProces(client proto.ClientClient, clientId int) {
	_, err := client.ClientJoin(context.Background(), &proto.ServerClientHasJoined{t : time, message : "Participant " + strconv.Itoa(clientId) + " joined Chitty-Chat at Lamport time " + strconv.Itoa(time) })

	if err != nil {
        log.Fatalf("Not working 4")
    }
}






