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
	myAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP("udp", myAddress)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func listenTerminal() {
	for {
		// TODO
	}
}

func listenOtherProcesses(connection *net.UDPConn) {
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
	go listenOtherProcesses(connection)
}
