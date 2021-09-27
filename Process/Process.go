package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SafeInt type used both to store clock and the number of replies
type SafeInt struct {
	mutex sync.Mutex
	value int
}

func (si *SafeInt) increment() (ret int) {
	si.mutex.Lock()
	si.value++
	ret = si.value
	si.mutex.Unlock()
	return ret
}

func (si *SafeInt) toZero() {
	si.mutex.Lock()
	si.value = 0
	si.mutex.Unlock()
}

func (clk *SafeInt) next(otherValue int) {
	defer updatePrompt()

	clk.mutex.Lock()
	defer clk.mutex.Unlock()

	if clk.value > otherValue {
		clk.value++
		return
	}

	clk.value = otherValue + 1
	return
}

// Possbible processes states
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

func (stt *SafeState) changeTo(value state) {
	stt.mutex.Lock()
	stt.value = value
	stt.mutex.Unlock()
}

// processMessage model the messages exchanged between processes
type processMessage struct {
	Id         int
	ClockValue int
	Text       string
}

// Global variables
var (
	myId               int
	myClock            SafeInt
	myState            SafeState = SafeState{value: RELEASED}
	myQueue            []*net.UDPConn
	numProcesses       int
	serverConn         *net.UDPConn
	clientConn         []*net.UDPConn
	sharedResourceConn *net.UDPConn
	clkValueAtReqst    int
	numReplies         SafeInt
	err                error
	printIsOk          = make(chan bool)
	enoughReplies      = make(chan bool)
)

const (
	SharedResourcePort = ":10001"
	sleepDuration      = 5 * time.Second
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func newListenConn(port string) *net.UDPConn {
	localAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	checkError(err)

	connection, err := net.ListenUDP("udp", localAddress)
	checkError(err)

	return connection
}

func newDialConn(port string) *net.UDPConn {
	remoteAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	checkError(err)

	connection, err := net.DialUDP("udp", nil, remoteAddress)
	checkError(err)

	return connection
}

func sendMessageTo(text string, connection *net.UDPConn) {
	message := processMessage{Id: myId, ClockValue: myClock.value, Text: text}
	bytes, err := json.Marshal(message)
	checkError(err)

	_, err = connection.Write(bytes)
	checkError(err)
}

func requestCriticalSection() {
	clkValueAtReqst = myClock.increment()
	printIsOk <- true

	numReplies.toZero()
	for i, connection := range clientConn {
		if i != myId-1 {
			sendMessageTo("REQUEST", connection)
		}
	}
	<-enoughReplies
}

func useCriticalSection() {
	myState.changeTo(HELD)
	fmt.Print("\nEntrei na CS")
	updatePrompt()

	sendMessageTo("Oi CS ^-^", sharedResourceConn)
	time.Sleep(sleepDuration)
}

func releaseCriticalSection() {
	fmt.Print("\nSai da CS")
	myState.changeTo(RELEASED)
	updatePrompt()

	for _, connection := range myQueue {
		sendMessageTo("REPLY", connection)
	}
	myQueue = []*net.UDPConn{}
}

func tryEnterCriticalSection() {
	switch myState.value {
	case HELD:
		fmt.Print("x ignorado")
		printIsOk <- true

	case WANTED:
		printIsOk <- true

	case RELEASED:
		myState.changeTo(WANTED)

		requestCriticalSection()
		useCriticalSection()
		releaseCriticalSection()

	default:
		panic(errors.New("Unexpected state"))
	}
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

func updatePrompt() {
	fmt.Printf("\nProcess %v, Clock %v, %v> ",
		myId, myClock.value, myState.value)
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
}

func resolveMessage(message processMessage) {
	switch message.Text {
	case "REQUEST":
		senderConn := clientConn[message.Id-1]

		if myState.value == RELEASED {
			myClock.next(message.ClockValue)
			sendMessageTo("REPLY", senderConn)

		} else if myState.value == HELD ||
			(clkValueAtReqst < message.ClockValue ||
				(clkValueAtReqst == message.ClockValue && myId < message.Id)) {

			myClock.next(message.ClockValue)
			myQueue = append(myQueue, senderConn)
		}

	case "REPLY":
		myClock.next(message.ClockValue)

		if numReplies.increment() == numProcesses-1 {
			enoughReplies <- true
		}
	}
}

func listenOtherProcesses() {
	jsonDecoder := json.NewDecoder(serverConn)

	for {
		var message processMessage
		jsonDecoder.Decode(&message)
		resolveMessage(message)
	}
}

func main() {
	myId, err = strconv.Atoi(os.Args[1])
	checkError(err)

	processesPorts := os.Args[2:]

	myPort := processesPorts[myId-1]
	serverConn = newListenConn(myPort)
	defer serverConn.Close()

	numProcesses = len(processesPorts)
	for i, port := range processesPorts {
		clientConn = append(clientConn, newDialConn(port))
		defer clientConn[i].Close()
	}

	sharedResourceConn = newDialConn(SharedResourcePort)
	defer sharedResourceConn.Close()

	go listenOtherProcesses()
	listenTerminal()
}
