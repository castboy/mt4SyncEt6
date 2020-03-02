package mt4SyncEt6_test

import (
	"fmt"
	zmqPKG "mt4SyncEt6/zmq"
	"mt4SyncEt6/zmq/dealer_rep"
	"runtime"
	"testing"
)

func TestZMQ(t *testing.T) {
	addr := "120.131.3.46:5590"
	zm := zmqPKG.ZMQ4Work(addr)
	//request := "{\"__api\":\"AccountInfo\",\"login\":2500659111,\"trade_system_type\":1}"
	request := "{\"__api\":\"AccountInfo\",\"login\":2500693,\"lang\":\"cn\",\"php_time\":1579155558,\"php_host\":\"dev-global-demo-08\",\"trade_system_type\":\"1\"}"
	reply, err := zm.Request(request)
	if err != nil {
		panic(err)
	}
	fmt.Println("reply======", reply)
}

/*func TestGoZMQ(t *testing.T) {
	addrRouter := "tcp://localhost:5559"
	addrDealer := "120.131.3.46:5590"
	gozmq := zmqPKG.GetGoZMQ()
	gozmq.RouterStart(addrRouter, addrDealer)
	gozmq.ClientStart(addrRouter)

	//msg req rep
	request := "{\"__api\":\"AccountInfo\",\"login\":2500693,\"lang\":\"cn\",\"php_time\":1579155558,\"php_host\":\"dev-global-demo-08\",\"trade_system_type\":\"1\"}"

	recv := gozmq.ZMQWork(request, 0)

	fmt.Printf("Received reply [%s]\n", recv)

}*/

func TestDealerRep(t *testing.T) {
	go dealer_rep.Dealer1_rep()
	go dealer_rep.Dealer2_rep()
	go dealer_rep.Dealer2_rep_server()

	for {
		runtime.GC()
	}
}

func TestDealerRep2(t *testing.T) {
	go dealer_rep.DealerRep11()
	go dealer_rep.DealerRep22()
	go dealer_rep.DealerRepServer()

	for {
		runtime.GC()
	}
}
