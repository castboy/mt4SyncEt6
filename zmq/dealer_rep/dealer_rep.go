package dealer_rep

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"time"
)

func Dealer1_rep() {
	//ZMQ ctx
	ctx, _ := zmq.NewContext()
	defer ctx.Term()
	//ZMQ socket
	zmqRequest, _ := ctx.NewSocket(zmq.DEALER)
	defer zmqRequest.Close()
	zmqRequest.Connect("tcp://localhost:5555")

	fmt.Println("Dealer1 is running")

	for i := 0; i < 10; i++ {
		zmqRequest.Send("", zmq.SNDMORE)
		zmqRequest.Send("Hello1", 0)

		reply1, _ := zmqRequest.Recv(0)
		if reply1 != "" {
		}
		reply, _ := zmqRequest.Recv(0)
		fmt.Printf("Received reply %d [%s]\n", i, reply)
	}
}

func Dealer2_rep() {
	//ZMQ ctx
	ctx, _ := zmq.NewContext()
	defer ctx.Term()

	//ZMQ socket
	zmqRequest, _ := ctx.NewSocket(zmq.DEALER)
	defer zmqRequest.Close()
	zmqRequest.Connect("tcp://localhost:5555")

	fmt.Println("Dealer2 is running")

	for i := 0; i < 10; i++ {
		zmqRequest.Send("", zmq.SNDMORE)
		zmqRequest.Send("Hello2", 0)

		reply1, _ := zmqRequest.Recv(0)
		if reply1 != "" {
		}
		reply, _ := zmqRequest.Recv(0)
		fmt.Printf("Received reply %d [%s]\n", i, reply)
	}
}

func Dealer2_rep_server() {
	//Context and sockect
	ctx, _ := zmq.NewContext()
	socket, _ := ctx.NewSocket(zmq.REP)
	defer ctx.Term()
	defer socket.Close()
	//Bind
	socket.Bind("tcp://*:5555")

	fmt.Println("Server is running")

	for {
		msg, _ := socket.Recv(0)
		fmt.Printf("Received :%s\n", msg)
		//Sleep
		time.Sleep(time.Second)
		//Send reply back to
		reply := fmt.Sprintf("World")
		socket.SendBytes([]byte(reply), 0)
	}
}

func DealerRep11() {
	context, _ := zmq.NewContext()
	defer context.Term()

	// Socket to talk to clients
	requester, _ := context.NewSocket(zmq.DEALER)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

	fmt.Print("DEALER1 is running")
	for i := 0; i < 10; i++ {
		//fmt.Print("req1 is running")
		//send
		requester.Send("", zmq.SNDMORE)
		requester.Send("Hello1", 0)

		//receive
		reply1, _ := requester.Recv(0)
		if reply1 != "" {

		}
		reply, _ := requester.Recv(0)
		fmt.Printf("Received reply %d [%s]\n", i, reply)
	}
}

func DealerRep22() {
	context, _ := zmq.NewContext()
	defer context.Term()

	// Socket to talk to clients
	requester, _ := context.NewSocket(zmq.DEALER)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

	fmt.Print("DEALER2 is running")
	for i := 0; i < 10; i++ {
		//fmt.Print("req1 is running")
		//send
		requester.Send("", zmq.SNDMORE)
		requester.Send("Hello2", 0)

		//receive
		reply1, _ := requester.Recv(0)
		if reply1 != "" {

		}

		reply, _ := requester.Recv(0)
		fmt.Printf("Received reply %d [%s]\n", i, reply)
	}
}

func DealerRepServer() {
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.REP)
	defer context.Term()
	defer socket.Close()
	socket.Bind("tcp://*:5555")

	fmt.Println("server is running ...")
	// Wait for messages
	for {
		msg, _ := socket.Recv(0)
		println("Received ", string(msg))
		// do some fake "work"
		time.Sleep(time.Second)
		// send reply back to client
		reply := fmt.Sprintf("World")
		socket.SendBytes([]byte(reply), 0)
	}
}
