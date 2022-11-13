# Ricart-Agrawala Algorithm

This is a simple implementation of the Ricart-Agrawala algorithm in Go. It is a distributed mutual exclusion algorithm that allows processes to request and release critical sections.

## Example - 3 processes


To run the program, you need to have Go installed. Then, open 4 terminals and run the following commands:

1st terminal:

    go run SharedResource/SharedResource.go

2nd terminal - Process 1:

    go run Process/Process.go 1 :1002 :1004 :1004

3rd terminal - Process 2:

    go run Process/Process.go 2 :1002 :1003 :1004

4th terminal - Process 3:
    
    go run Process/Process.go 3 :1002 :1003 :1004


This will start 3 processes and a shared resource. The first argument is the process ID, the second to fourth arguments are the port of the processes in order of their IDs. The shared resource will be listening on the fixed port 1001.

On each terminal running a process, it will show the process ID, its clock, and the current state. The state can be either "RELEASED", "WANTED", or "HELD". The clock is incremented every time the process requests or releases the critical section.

You can request the critical section by typing "x" on the terminal. The critical section will be released automatically after 5 seconds.

### Test case

The following diagram shows the execution of the test case. The process 1 requests the critical section at time 1. The process 3 requests the critical section at time 3. The process 2 requests the critical section at time 5. The process 1 releases the critical section at time 7. The process 3 releases the critical section at time 8. The process 2 releases the critical section at time 9.

![Legend](https://i.imgur.com/MyxKjwK.png)

![Test case](https://i.imgur.com/adJzosY.png)

![Terminals](https://i.imgur.com/qkk5xb4.png)

