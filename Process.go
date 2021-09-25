package main

import (
	"net"
	"os"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func createLocalConnection(port string) (*net.UDPConn, error) {
	// TODO
	return nil, nil
}

func listenTerminal() {
	for {
		// TODO
	}
}

func listenProcesses(connection *net.UDPConn) {
	for {
		// TODO
	}
}

func main() {
	myId, err := strconv.Atoi(os.Args[1])
	checkError(err)
	processesPorts := os.Args[2:]
	myPort := processesPorts[myId-1]
	connection, err := createLocalConnection(myPort)
	checkError(err)
	defer connection.Close()

	go listenTerminal()
	go listenProcesses(connection)
}
