package mt4SyncEt6

import (
	"encoding/json"
	"fmt"
	"mt4SyncEt6/session"
	"sort"
	"testing"
	"time"
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

func TestConvet(t *testing.T) {

	//prepare format and space for json convert
	var infoSlice = make([]SessionInfo, 300)
	for k, _ := range infoSlice {
		if k < 300 {
			infoSlice[k].Trade_session = make(map[string]string, 7)
		}

	}
	//json convert
	err := json.Unmarshal([]byte(session.SessionStr), &infoSlice)
	if err != nil {
		t.Fatalf("json.Unmarshal error:%+v\n", err)
	}

	//change mt4 to et6
	et6Sessions := make([]Et6Session, 0)
	for _, v := range infoSlice {
		if v.Symbol_name == "" {
			t.Log("Empty symbol info")
			break
		}
		trade_session := session.Mt4ToEt6(v.Trade_session)
		et6Session := Et6Session{
			Symbol_name:   v.Symbol_name,
			Trade_session: trade_session,
		}
		et6Sessions = append(et6Sessions, et6Session)
	}
	t.Log(et6Sessions)
	// save to db
	for _, v := range et6Sessions {
		SessionToDB(v)
	}
}

func Test_Sessiion(t *testing.T){
	sesses:=[]*SessionNew{}
	engine,_:=NewKiteEngineXorm()
	engine.Table("session").Find(&sesses)

}




type SessionDst struct {
	ID       int          `xorm:"id autoincr"`
	SourceID int          `xorm:"source_id"`
	Type     SessionType  `xorm:"type"`
	DstType     DSTType   `xorm:"dst_type"`
	TimeSpan string       `xorm:"time_span"`
}


var quote = map[int]map[time.Weekday][]string{}
var trade = map[int]map[time.Weekday][]string{}

func Test_DecodeSession(t *testing.T) {
	insertData()
	//decodeData()
}


func insertData() {
	xlive, err := NewProduceEngineXorm()
	//x, err := xorm.NewEngine("mysql", "root:wang1234@/trading_system?charset=utf8")
	if err != nil {
		panic(err)
	}
	x,err:=NewLocalhostEngineXorm()
	//x, err := xorm.NewEngine("mysql", "root:wang1234@/trading_system?charset=utf8")
	if err != nil {
		panic(err)
	}

	ss := []Session{}

	err = xlive.Table("session").Find(&ss)
	if err != nil {
		panic(err)
	}

	for i := range ss {
		if ss[i].Type == Quote {
			if quote[ss[i].SourceID] == nil {
				quote[ss[i].SourceID] = map[time.Weekday][]string{}
			}

			quote[ss[i].SourceID][ss[i].Weekday] = append(quote[ss[i].SourceID][ss[i].Weekday], ss[i].TimeSpan)
			continue
		}

		if trade[ss[i].SourceID] == nil {
			trade[ss[i].SourceID] = map[time.Weekday][]string{}
		}

		trade[ss[i].SourceID][ss[i].Weekday] = append(trade[ss[i].SourceID][ss[i].Weekday], ss[i].TimeSpan)
	}

	var sort = func(s map[int]map[time.Weekday][]string) {
		for i := range s {
			for j := range s[i] {
				sort.Strings(s[i][j])
			}
		}
	}

	sort(quote)
	sort(trade)

	var fillZero = func(s map[int]map[time.Weekday][]string) {
		for i := range s {
			for j := range s[i] {
				if len(s[i][j]) == 3 {
					continue
				}
				if len(s[i][j]) == 2 {
					s[i][j] = append(s[i][j], "00:00-00:00")
				}
				if len(s[i][j]) == 1 {
					s[i][j] = append(s[i][j], "00:00-00:00", "00:00-00:00")
				}
			}

			if s[i][time.Saturday] == nil {
				s[i][time.Saturday] = []string{"00:00-00:00", "00:00-00:00", "00:00-00:00"}
			}

			if s[i][time.Sunday] == nil {
				s[i][time.Sunday] = []string{"00:00-00:00", "00:00-00:00", "00:00-00:00"}
			}
		}
	}

	fillZero(quote)
	fillZero(trade)

	x, _ = NewKiteEngineXorm()

	var insert = func(s map[int]map[time.Weekday][]string, sessionType SessionType) {
		L := len(s)
		for i := 1; i < L + 1 ; i++ {

			q, err := json.Marshal(s[i])
			if err != nil {
				panic(err)
			}

			dst := &SessionDst{
				SourceID: i,
				Type:sessionType,
				DstType:NonDST,
				TimeSpan: string(q),
			}

			if len(q) == 0 {
				continue
			}
			if i==27{
				fmt.Println("HK50",dst)
				continue
			}
			_, err = x.Table("session").Insert(dst)
			if err != nil {
				panic(err)
			}
		}
	}

	insert(quote, Quote)
	insert(trade, Trade)
}
