package main

import (
	"net"
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

func listenProcesses(connection *net.UDPConn) {
	for {
		// TODO
	}
}

func main() {
	myPort := ":10001"
	connection, err := createLocalConnection(myPort)
	checkError(err)
	defer connection.Close()

	listenProcesses(connection)
}
