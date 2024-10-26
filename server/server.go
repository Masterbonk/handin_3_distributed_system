package main

import (
	cc "ChittyChat/ChittyChat"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var clientId int = 0

type server struct {
	cc.UnimplementedChittyChatServer
	clients map[int]cc.ChittyChat_ChatServer

	lamport int32

	mu sync.Mutex // used to ensure only one goroutine can write at a time
}

// MessageType 0 is client joined/left
// MessageType 1 is client post
func (s *server) broadcast(lam int32, msg string, clientName string, messageType int32) {
	addToLamport(lam, &s.lamport)

	if len(msg) > 128{
		msg = msg[1:128]
	}

	for _, ss := range s.getClients() {

		if err := ss.Send(&cc.ServerMessage{Lamport: s.lamport, Msg: msg, ClientName: clientName, MessageType: messageType}); err != nil {
			log.Printf("broadcast err %s: %v", clientName, err)
		} 
	}


	log.Printf("Time: %d, %s: %s\n",s.lamport, clientName, msg)
	
}

func addToLamport(inputLamport int32, ourLamport *int32) {
	if inputLamport > *ourLamport {
		*ourLamport = inputLamport
	}
	*ourLamport++
}

func (s *server) broadcastMessage(in cc.ClientMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.broadcast(in.Lamport, in.Msg, in.ClientName, 1)
}

func (s *server) addClient(uid int, in cc.ClientMessage, srv cc.ChittyChat_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[uid] = srv
	var msg string = fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", in.ClientName, s.lamport)
	s.broadcast(s.lamport, msg, in.ClientName, 0)
}

func (s *server) removeClient(uid int, username string, lamport int32, srv cc.ChittyChat_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// maybe close srv???

	delete(s.clients, uid)

	var msg string = fmt.Sprintf("Participant %s left Chitty-Chat at Lamport time %d", username, s.lamport)
	s.broadcast(s.lamport, msg, username, 0)
}

func (s *server) Chat(srv cc.ChittyChat_ChatServer) error {
	in, err := srv.Recv()
	s.addClient(clientId, *in, srv)

	userId := clientId
	username := in.ClientName

	clientId++

	for {
		in, err = srv.Recv()

		// if reached end of input
		if err == io.EOF {
			break
		} else {
			addToLamport(in.Lamport, &s.lamport)
		}
		// if encountered another error
		if err != nil {
			return err
		}

		s.broadcastMessage(*in)
	}

	s.removeClient(userId, username, s.lamport, srv)

	return nil
}

func newServer() *server {
	s := &server{clients: make(map[int]cc.ChittyChat_ChatServer)}
	s.lamport = 0
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
	ip := "localhost"
	//ip := "192.168.43.80"

	// listen to port
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",ip, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("Now listening to port: %d", port)
	}

	grpcServer := grpc.NewServer()
	cc.RegisterChittyChatServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
