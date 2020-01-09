package mt4SyncEt6

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"mt4SyncEt6/decimal"
	"strconv"
	"time"
)

func GetEngine() *xorm.Engine {
	enginEt6, _ := NewET6EngineXorm()
	return enginEt6
}

func SessionToDB(v Et6Session) {
	enginEt6, err := NewET6EngineXorm()
	if err != nil {
		return
	}
	sourceItem := new(Source)
	//_, err = enginEt6.Table("source").Where("source=?", v.Symbol_name).NoAutoCondition(true).Get(sourceItem)
	hit, err := enginEt6.Table("source").Where("source=?", v.Symbol_name).NoAutoCondition(true).Get(sourceItem)
	if !hit && err == nil {
		return
	}
	//insert session to db
	// map
	for k, tradeSession := range v.Trade_session {
		//slice
		for _, span := range tradeSession {
			sess := Session{}
			sess.Type = 0
			sess.SourceID = sourceItem.ID
			weekDay, _ := strconv.Atoi(k)
			if weekDay == 0 {
				sess.Weekday = time.Sunday
			} else if weekDay == 1 {
				sess.Weekday = time.Monday
			} else if weekDay == 2 {
				sess.Weekday = time.Tuesday
			} else if weekDay == 3 {
				sess.Weekday = time.Wednesday
			} else if weekDay == 4 {
				sess.Weekday = time.Thursday
			} else if weekDay == 5 {
				sess.Weekday = time.Friday
			} else if weekDay == 6 {
				sess.Weekday = time.Saturday
			}
			sess.TimeSpan = span

			//trade type
			sess.Type = 0
			_, err = enginEt6.Table("session").Omit("id").InsertOne(sess)
			sess.Type = 1
			_, err = enginEt6.Table("session").Omit("id").InsertOne(sess)
			if err != nil {
				fmt.Println("err")
			}
		}
	}
}

func NewET6EngineXorm() (*xorm.Engine, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"glmt4dev_wr", "mt4geed0Uokohphai1UNgeep5ae", "devcondb.r62g.cn",
		"3306", "trading_system")

	mt4Engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	mt4Engine.SetMaxOpenConns(100)
	mt4Engine.SetMaxIdleConns(20)
	mt4Engine.SetConnMaxLifetime(1800 * time.Second)
	//engine.DatabaseTZ = time.UTC
	//engine.TZLocation = time.UTC

	return mt4Engine, nil
}

type SwapInfo struct {
	Symbol    string  `json:"symbol"`
	SwapLong  float64 `json:"swap_long"`
	SwapShort float64 `json:"swap_short"`
	Swap3Day  string  `json:"Swap3Days"`
	SourceCN  string  `json:"symbol_cn"`
}

type SessionInfo struct {
	Symbol_name   string            `json:"symbol_name"`
	Trade_session map[string]string `json:"trade_session"`
}

type Et6Session struct {
	Symbol_name   string              `json:"symbol_name"`
	Trade_session map[string][]string `json:"trade_session"`
}

type Source struct {
	ID             int             `json:"id" xorm:"id autoincr"`
	Source         string          `json:"source" xorm:"source"`
	SourceCN       string          `json:"source_cn" xorm:"source_cn"`
	SourceType     SourceType      `json:"source_type" xorm:"source_type"`
	Digits         int             `json:"digits" xorm:"digits"`
	Multiply       decimal.Decimal `json:"multiply" xorm:"multiply"`
	ContractSize   decimal.Decimal `json:"contract_size" xorm:"contract_size"`
	StopsLevel     int             `json:"stops_level" xorm:"stops_level"`
	ProfitMode     ProfitMode      `json:"profit_mode" xorm:"profit_mode"`
	ProfitCurrency string          `json:"profit_currency" xorm:"profit_currency"`
	MarginMode     MarginMode      `json:"margin_mode" xorm:"margin_mode"`
	MarginCurrency string          `json:"margin_currency" xorm:"margin_currency"`
	SwapType       SwapType        `json:"swap_type" xorm:"swap_type"`
	SwapLong       decimal.Decimal `json:"swap_long" xorm:"swap_long"`
	SwapShort      decimal.Decimal `json:"swap_short" xorm:"swap_short"`
	SwapCurrency   string          `josn:"swap_currency" xorm:"swap_currency"`
	Swap3Day       time.Weekday    `json:"swap_3_day" xorm:"swap_3_day"`
}

type Session struct {
	ID       int          `xorm:"id autoincr"`
	SourceID int          `xorm:"source_id"`
	Type     SessionType  `xorm:"type"`
	Weekday  time.Weekday `xorm:"weekday"`
	TimeSpan string       `xorm:"time_span"`
}

type AccountGroup struct {
	ID              int             `xorm:"id autoincr"`
	Name            string          `xorm:"name"`
	DepositCurrency string          `xorm:"deposit_currency"`
	MarginStopOut   decimal.Decimal `xorm:"margin_stop_out"`
	MarginMode      MarginCalcMode  `xorm:"margin_mode"`
}

type ConGroupSec struct {
	ID             int             `xorm:"id autoincr"`
	GroupId        int             `xorm:"group_id"`
	SecurityId     int             `xorm:"security_id"`
	EnableSecurity int             `xorm:"enable_security"`
	EnableTrade    int             `xorm:"enable_trade"`
	LotMin         decimal.Decimal `xorm:"lot_min"`
	LotMax         decimal.Decimal `xorm:"lot_max"`
	LotStep        decimal.Decimal `xorm:"lot_step"`
	SpreadDiff     int             `xorm:"spread_diff"`
	Commission     decimal.Decimal `xorm:"commission"`
}

type MarginCalcMode int

const (
	DoubleLegs MarginCalcMode = iota
	LargerLeg
	NetLeg
)

type SessionType int

const (
	Quote SessionType = iota
	Trade
)

type (
	ProfitMode int
	SwapType   int
	MarginMode int
	SourceType int
)

const (
	MarginForex MarginMode = iota
	MarginCfd
	MarginFutures
	MarginCfdIndex
	MarginCfdLeverage
)

const (
	ByPoints SwapType = iota
	ByMoney
	ByInterest
	ByMoneyInMarginCurrency
	ByInterestOfCfds
	ByInterestOfFutures
)

const (
	SourceFx SourceType = iota
	SourceMetal
	SourceEnergy
	SourceIndex
	SourceCrypto
)

const (
	ProfitForex ProfitMode = iota
	ProfitCfd
	ProfitFutures
)
