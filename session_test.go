package mt4SyncEt6

import (
	"encoding/json"
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

func TestConvet(t *testing.T)  {

	//prepare format and space for json convert
	var infoSlice=make([]SessionInfo,300)
	for k,_:=range infoSlice{
		if k<300{
			infoSlice[k].Trade_session= make(map[string]string,7)
		}

	}
	//json convert
	err:=json.Unmarshal([]byte(SesssionStr),&infoSlice)
	if err!=nil{
		t.Fatalf("json.Unmarshal error:%+v\n",err)
	}

	//change mt4 to et6
	et6Sessions:=make([]Et6Session,0)
	for _,v:=range infoSlice{
		if v.Symbol_name==""{
			t.Log("Empty symbol info")
			break
		}
		trade_session:=session.Mt4ToEt6(v.Trade_session)
		et6Session:=Et6Session{
			Symbol_name:v.Symbol_name,
			Trade_session:trade_session,
		}
		et6Sessions = append(et6Sessions, et6Session)
	}
	t.Log(et6Sessions)
	// save to db
	for _,v:=range et6Sessions{
		SessionToDB(v)
	}
}


