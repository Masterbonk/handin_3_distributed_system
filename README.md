# Handin 3 distributed system
## How the program works
To test the program, one needs to start two or more terminals on the same computer simulating different computers. All terminals needs to be in the **handin_3_distributed_system** directory

The first terminal, which will be called serverTerm from now on, can start the server by using the following command in the terminal:
`go run server/server.go`

The terminal will write *Now listening to port: 50051* once the server is up and running.

Now one can connect with other terminals as clients. This can be done by writing the following command in the other terminal:
`go run client/client.go`
This will create a user called anouminous. One is able to add a flag to specify ones name, using '-u name', so the following command will start a client with the name Tom:
`go run client/client.go -u Tom`

Everytime a new client joins, all clients will get a message from the server stating their name. 

To send messages, the clients can be written in, where using *enter* will send the current message. If the message is the word *shutdown*, the specific client will send a leave message and stop itself.