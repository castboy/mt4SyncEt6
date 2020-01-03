package mt4SyncEt6

import (
	"fmt"
	"mt4SyncEt6/session"
	"testing"
)

func TestMt4ToEt6(t *testing.T) {
	mt4 := map[string]string{
		"0": "0:0-0:3,3:0-5:0,23:0-24:0",
		"1": "0:0-12:0,13:0-24:0,0:0-0:0",
		"2": "9:0-9:5,12:0-13:0,23:0-24:0",
		"3": "12:0-20:0,20:10-21:3,22:0-24:0",
	}

	et6 := session.Mt4ToEt6(mt4)
	fmt.Println(et6)
}
