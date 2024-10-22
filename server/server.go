package main

import (
	cc "ChittyChat/ChittyChat"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

var clientId int = 0

type server struct {
	cc.UnimplementedChittyChatServer
	clients map[string]cc.ChittyChat_ChatServer

	mu sync.Mutex // used to ensure only one goroutine can write at a time
}

// MessageType 0 is client joined/left
// MessageType 1 is client post
func (s *server) broadcast(lamport int32, msg string, clientName string, messageType int32) {

	for _, ss := range s.getClients() {
		if err := ss.Send(&cc.ServerMessage{Lamport: lamport, Msg: msg, ClientName: clientName, MessageType: messageType}); err != nil {
			log.Printf("broadcast err: %v", err)
		}
	}
}

func (s *server) broadcastMessage(in cc.ClientMessage) {
	s.broadcast(in.Lamport, in.Msg, in.ClientName, 1)
}

func (s *server) addClient(uid int, in cc.ClientMessage, srv cc.ChittyChat_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	index := strconv.Itoa(uid)
	s.clients[index] = srv

	var msg string = fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", in.ClientName, in.Lamport)
	s.broadcast(0, msg, in.ClientName, 0)
}

func (s *server) Chat(srv cc.ChittyChat_ChatServer) error {
	in, err := srv.Recv()
	s.addClient(clientId, *in, srv)
	clientId++

	for {
		in, err = srv.Recv()

		// if reached end of input
		if err == io.EOF {
			return nil
		}
		// if encountered another error
		if err != nil {
			return err
		}

		s.broadcastMessage(*in)
	}
}

func newServer() *server {
	s := &server{clients: make(map[string]cc.ChittyChat_ChatServer)}
	return s
}

func (s *server) getClients() []cc.ChittyChat_ChatServer {
	var cs []cc.ChittyChat_ChatServer

	for _, c := range s.clients {
		cs = append(cs, c)
	}
	return cs
}

func main() {
	port := 50051

	// listen to port
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Printf("Now listening to port: %d", port)
	}

	grpcServer := grpc.NewServer()
	cc.RegisterChittyChatServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
