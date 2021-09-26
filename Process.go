package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type SafeClock struct {
	mutex sync.Mutex
	value int
}

var (
	myId  int
	clock SafeClock
	done  = make(chan bool)
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

func tryEnterCriticalSection() {
	// TODO
	fmt.Println("Tentando entrar na região crítica...")
}

func incrementClock() {
	clock.mutex.Lock()
	clock.value++
	clock.mutex.Unlock()
}

func useInput(input string) {
	if strings.ToLower(input) == "x" {
		tryEnterCriticalSection()
		return
	}

	num, err := strconv.Atoi(input)
	if err == nil && num == myId {
		incrementClock()
	}
}

func listenTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Process %v, Clock %v> ", myId, clock.value)
		input, err := reader.ReadString('\n')
		checkError(err)
		input = input[:len(input)-1]
		useInput(input)
	}
	done <- true
}

func listenOtherProcesses(connection *net.UDPConn) {
	for {
		// TODO
	}
}

func main() {
	var err error
	myId, err = strconv.Atoi(os.Args[1])
	checkError(err)
	processesPorts := os.Args[2:]
	myPort := processesPorts[myId-1]
	connection, err := createLocalConnection(myPort)
	checkError(err)
	defer connection.Close()

	go listenTerminal()
	go listenOtherProcesses(connection)
	<-done
}
