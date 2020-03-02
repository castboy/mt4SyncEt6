package zmq

/*
import (
	zmq "github.com/alecthomas/gozmq"
)

type GoZmq struct {
	Context zmq.Context
	Socket  *zmq.Socket
}

var Gozmq *GoZmq

func GetGoZMQ() *GoZmq {
	if Gozmq == nil {
		Gozmq = new(GoZmq)
	}
	return Gozmq
}

func (mz *GoZmq) RouterStart(addrRouter string, addrDealer string) string {
	context, _ := zmq.NewContext()
	defer context.Close()

	frontend, _ := context.NewSocket(zmq.ROUTER)
	backend, _ := context.NewSocket(zmq.DEALER)
	defer frontend.Close()
	defer backend.Close()
	frontend.Bind("tcp://" + addrRouter)
	backend.Bind("tcp://" + addrDealer)

	// Initialize poll set
	toPoll := zmq.PollItems{
		zmq.PollItem{Socket: frontend, Events: zmq.POLLIN},
		zmq.PollItem{Socket: backend, Events: zmq.POLLIN},
	}
	for {
		_, _ = zmq.Poll(toPoll, -1)

		switch {
		case toPoll[0].REvents&zmq.POLLIN != 0:
			parts, _ := frontend.RecvMultipart(0)
			backend.SendMultipart(parts, 0)

		case toPoll[1].REvents&zmq.POLLIN != 0:
			parts, _ := backend.RecvMultipart(0)
			frontend.SendMultipart(parts, 0)
		}
	}
}

func (mz *GoZmq) ClientStart(addrRouter string) {
	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to talk to clients
	requester, _ := context.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://" + addrRouter)
	mz.Socket = requester
}

func (mz *GoZmq) ZMQWork(request string, tag zmq.SendRecvOption) string {
	mz.Socket.Send([]byte(request), tag)
	recv, _ := mz.Socket.Recv(tag)
	return string(recv)
}
*/
