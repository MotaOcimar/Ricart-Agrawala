package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type processMessage struct {
	Id    int
	Clock int
	Text  string
}

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
	jsonDecoder := json.NewDecoder(connection)
	for {
		var message processMessage
		jsonDecoder.Decode(&message)
		fmt.Printf("Process of id %v and clock %v says: \"%v\"\n",
			message.Id, message.Clock, message.Text)
	}
}

func main() {
	const myPort = ":10001"
	connection, err := createLocalConnection(myPort)
	checkError(err)
	defer connection.Close()

	listenProcesses(connection)
}
