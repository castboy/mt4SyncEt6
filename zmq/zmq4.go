package zmq

import (
	"fmt"
	zmq4 "github.com/pebbe/zmq4"
	"math/rand"
	"time"
)

type MessageZmq struct {
	context *zmq4.Socket
}

func (mz *MessageZmq) set_ID() string {
	flag := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	mz.context.SetIdentity(flag)
	return flag
}

func (mz *MessageZmq) Request(str string) (string, error) {
	identity := mz.set_ID()
	countI, err := mz.context.Send(identity, zmq4.SNDMORE)
	if err != nil {
		panic(err)
	}
	fmt.Println("send identity", countI)
	countStr, err := mz.context.Send(str, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("send countStr", countStr)
	return mz.context.Recv(0)
}

func ZMQ4Work(addr string) *MessageZmq {
	zmqContext, _ := zmq4.NewSocket(zmq4.DEALER)
	zmqContext.SetRcvtimeo(100 * time.Second)
	zmqContext.Connect("tcp://" + addr)
	serv := &MessageZmq{zmqContext}
	return serv
}
