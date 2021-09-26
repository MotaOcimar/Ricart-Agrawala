package main

import (
	"bufio"
	"errors"
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

type state string

const (
	RELEASED state = "RELEASED"
	WANTED   state = "WANTED"
	HELD     state = "HELD"
)

type SafeState struct {
	mutex sync.Mutex
	value state
}

var (
	myId      int
	myClock   SafeClock
	myState   SafeState = SafeState{value: RELEASED}
	done                = make(chan bool)
	printIsOk           = make(chan bool)
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

func requestCriticalSection() {
	// TODO:
}

func useCriticalSection() {
	// TODO:
}

func releaseCriticalSection() {
	// TODO:
}

func (stt *SafeState) changeTo(value state) {
	stt.mutex.Lock()
	stt.value = value
	stt.mutex.Unlock()
}

func tryEnterCriticalSection() {
	switch myState.value {
	case HELD:
		fmt.Println("x ignorado")
		printIsOk <- true

	case WANTED:
		printIsOk <- true

	case RELEASED:
		myState.changeTo(WANTED)
		printIsOk <- true

		requestCriticalSection()
		useCriticalSection()
		releaseCriticalSection()

	default:
		panic(errors.New("Unexpected state"))
	}
}

func (clk *SafeClock) increment() {
	clk.mutex.Lock()
	clk.value++
	clk.mutex.Unlock()
}

func useInput(input string) {
	if strings.ToLower(input) == "x" {
		go tryEnterCriticalSection()
		<-printIsOk
		return
	}

	num, err := strconv.Atoi(input)
	if err == nil && num == myId {
		myClock.increment()
	}
}

func listenTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Process %v, Clock %v, %v> ",
			myId, myClock.value, myState.value)
		input, err := reader.ReadString('\n')
		checkError(err)
		input = input[:len(input)-1]
		useInput(input)
	}
	done <- true
}

func listenOtherProcesses(connection *net.UDPConn) {
	for {
		// TODO:
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
