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

type SafeInt struct {
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

type processMessage struct {
	Id         int
	ClockValue int
	Text       string
}

var (
	myId            int
	myClock         SafeInt
	myState         SafeState = SafeState{value: RELEASED}
	myPort          string
	myConnection    *net.UDPConn // TODO: Fazer uma conexão para responder cada cliente e uma só para escutar
	myQueue         []string
	clkValueAtReqst int
	numReplies      SafeInt
	processesPorts  []string
	done            = make(chan bool)
	printIsOk       = make(chan bool)
	enoughReplies   = make(chan bool)
)

const (
	SharedResoursePort = ":10001"
	sleepDuration      = 5 * time.Second
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

func sendMessageTo(text string, port string) {
	receiverAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	checkError(err)

	message := processMessage{Id: myId, ClockValue: myClock.value, Text: text}
	bytes, err := json.Marshal(message)
	checkError(err)
	_, err = myConnection.WriteToUDP(bytes, receiverAddress)
	checkError(err)
}

func requestCriticalSection() {
	clkValueAtReqst = myClock.increment()
	printIsOk <- true

	numReplies.toZero()
	for _, port := range processesPorts {
		if port != myPort {
			sendMessageTo("REQUEST", port)
		}
	}
	<-enoughReplies
}

func useCriticalSection() {
	myState.changeTo(HELD)
	fmt.Print("\nEntrei na CS")
	updatePrompt()

	sendMessageTo("Oi CS ^-^", SharedResoursePort)
	time.Sleep(sleepDuration)
}

func releaseCriticalSection() {
	fmt.Print("\nSai da CS")
	myState.changeTo(RELEASED)
	updatePrompt()

	for _, port := range myQueue {
		sendMessageTo("REPLY", port)
	}
	myQueue = []string{}
}

func (stt *SafeState) changeTo(value state) {
	stt.mutex.Lock()
	stt.value = value
	stt.mutex.Unlock()
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
	done <- true
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

func resolveMessage(message processMessage) {
	switch message.Text {
	case "REQUEST":
		senderPort := processesPorts[message.Id-1]

		if myState.value == RELEASED {
			myClock.next(message.ClockValue)
			sendMessageTo("REPLY", senderPort)

		} else if myState.value == HELD ||
			(clkValueAtReqst < message.ClockValue ||
				(clkValueAtReqst == message.ClockValue && myId < message.Id)) {

			myClock.next(message.ClockValue)
			myQueue = append(myQueue, senderPort)
		}

	case "REPLY":
		myClock.next(message.ClockValue)

		numProcesses := len(processesPorts)
		if numReplies.increment() == numProcesses-1 {
			enoughReplies <- true
		}
	}
}

func listenOtherProcesses() {
	jsonDecoder := json.NewDecoder(myConnection)

	for {
		var message processMessage
		jsonDecoder.Decode(&message)
		resolveMessage(message)
	}
}

func main() {
	var err error
	myId, err = strconv.Atoi(os.Args[1])
	checkError(err)
	processesPorts = os.Args[2:]
	myPort = processesPorts[myId-1]
	myConnection, err = createLocalConnection(myPort)
	checkError(err)
	defer myConnection.Close()

	go listenTerminal()
	go listenOtherProcesses()
	<-done
}
