package mt4SyncEt6

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"mt4SyncEt6/decimal"
	"strconv"
	"strings"
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func GetEngine() *xorm.Engine {
	enginEt6, _ := NewET6EngineXorm()
	return enginEt6
}

/*func TimeConv(in int) (out string) {
	if in/60 == 0 && in%60 == 0 {
		out = strconv.Itoa(in/60) + ":" + "00" + ":" + "00"
		return
	}
	out = strconv.Itoa(in/60) + ":" + strconv.Itoa(in%60) + ":" + "59"
	return
}*/
func ConvToItemAddOneDay(in time.Time) (tf string, tt string, datef string) {
	d, _ := time.ParseDuration("24h")
	in = in.Add(d)
	inStr := in.Format(TIME_LAYOUT)
	fmt.Println("inStr==========", inStr)
	dataSlice := strings.Split(string(inStr), " ")
	datef = dataSlice[0]
	tf = "00" + ":" + "00" + ":" + "00"
	if in.Hour() == 0 && in.Minute() == 0 {
		timeInSlice := strings.Split(string(dataSlice[1]), ":")
		tt = timeInSlice[0] + ":" + timeInSlice[1] + ":" + "00"
		return
	}
	timeOutSlice := strings.Split(string(dataSlice[1]), ":")
	tt = timeOutSlice[0] + ":" + timeOutSlice[1] + ":" + "59"
	return
}
func ConvertInsameSay(in time.Time) (out, timeDate string) {
	inStr := in.Format(TIME_LAYOUT)
	dataSlice := strings.Split(string(inStr), " ")
	timeDate = dataSlice[0]
	out = dataSlice[1]
	return
}
func ConvFromItem(in time.Time) (ff string, ft string, datef string) {
	inStr := in.Format(TIME_LAYOUT)
	dataSlice := strings.Split(string(inStr), " ")
	datef = dataSlice[0]
	if in.Hour() == 0 && in.Minute() == 0 {
		timeInSlice := strings.Split(string(dataSlice[1]), ":")
		ff = timeInSlice[0] + ":" + timeInSlice[1] + ":" + "00"
		ft = "00" + ":" + "00" + ":" + "00"
		return
	}
	timeOutSlice := strings.Split(string(dataSlice[1]), ":")
	ff = timeOutSlice[0] + ":" + timeOutSlice[1] + ":" + "59"
	ft = "00" + ":" + "00" + ":" + "00"
	return
}

func ConvToItem(in time.Time) (tf string, tt string, datef string) {
	inStr := in.Format(TIME_LAYOUT)
	dataSlice := strings.Split(string(inStr), " ")
	datef = dataSlice[0]
	tf = "00" + ":" + "00" + ":" + "00"
	if in.Hour() == 0 && in.Minute() == 0 {
		timeInSlice := strings.Split(string(dataSlice[1]), ":")
		tt = timeInSlice[0] + ":" + timeInSlice[1] + ":" + "00"
		return
	}
	timeOutSlice := strings.Split(string(dataSlice[1]), ":")
	tt = timeOutSlice[0] + ":" + timeOutSlice[1] + ":" + "59"
	return
}
func TimeConv(timeDate string, in int) (strTime string) {
	var hour string
	var minute string
	if in/60 < 10 {
		hour = "0" + strconv.Itoa(in/60)
	} else {
		hour = strconv.Itoa(in / 60)
	}
	if in%60 < 10 {
		minute = "0" + strconv.Itoa(in%60)
	} else {
		minute = strconv.Itoa(in % 60)
	}
	strTime = timeDate + " " + hour + ":" + minute + ":" + "00"
	fmt.Println("TimeConv====", strTime)
	return
}
func IsSameDay(in1 time.Time, in2 time.Time) (flag TimeDay, isSame bool) {
	if in1.Weekday() != in2.Weekday() {
		return DayDIff, false
	}
	if in1.Weekday() == in2.Weekday() && in1.Minute() == in2.Minute() && in1.Hour() == in2.Hour() {
		return DaySameTimeSame, false
	}
	return DaySameTimeDIff, true
}
func TimeConvToUTC(timeString string) time.Time {
	loc, _ := time.LoadLocation("Europe/Moscow")
	//loc, _ = time.LoadLocation("Local")
	tMosc, _ := time.ParseInLocation(TIME_LAYOUT, timeString, loc)
	tUTC := tMosc.UTC()
	fmt.Println(tUTC.Format(TIME_LAYOUT))
	return tUTC
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
	/*dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
	"glmt4dev_wr", "mt4geed0Uokohphai1UNgeep5ae", "devcondb.r62g.cn",
	"3306", "trading_system")*/
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"root", "123456789", "localhost", "3306", "trading_system")
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

func NewProduceEngineXorm() (*xorm.Engine, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"et6_ro", "OhQu9Unei7iezair4Oe0", "rm-f2z35ztaojigwx291ro.mysql.eu-west-1.rds.aliyuncs.com",
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

func NewLocalhostEngineXorm() (*xorm.Engine, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"root", "123456789", "localhost",
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

type Holiday struct {
	ID          int             `json:"index" xorm:"-"`
	Enable      bool            `json:"enable" xorm:"enable"`
	Date        string          `json:"date" xorm:"date"`
	From        string          `json:"from" xorm:"from"`
	To          string          `json:"to" xorm:"to"`
	Category    HolidayCategory `json:"-" xorm:"category"`
	Symbol      string          `json:"symbol" xorm:"symbol"`
	Description string          `json:"description" xorm:"description"`
}
type HolidayMt4 struct {
	ID          int             `json:"index" xorm:"id"`
	Enable      bool            `json:"enable" xorm:"enable"`
	Date        string          `json:"date" xorm:"date"`
	From        int             `json:"from" xorm:"from"`
	To          int             `json:"to" xorm:"to"`
	Category    HolidayCategory `json:"-" xorm:"category"`
	Symbol      string          `json:"symbol" xorm:"symbol"`
	Description string          `json:"description" xorm:"description"`
}
type Security struct {
	ID           int      `json:"id" xorm:"id"`
	SecurityName string   `json:"security_name" xorm:"security_name"`
	Description  string   `json:"description" xorm:"description"`
	Symbols      []string `json:"symbols" xorm:"-"`
}
type SwapInfo struct {
	Symbol    string  `json:"symbol"`
	SwapLong  float64 `json:"swap_long"`
	SwapShort float64 `json:"swap_short"`
	Swap3Day  string  `json:"Swap3Days"`
	SourceCN  string  `json:"symbol_cn"`
}

type CurrencyInfo struct {
	Symbol   string `json:"symbol"`
	Currency string `json:"currency"`
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

// Symbol represents a instance of symbol
type Symbol struct {
	ID            int             `json:"id" xorm:"id autoincr"`
	Index         int             `json:"index" xorm:"index"`
	Symbol        string          `json:"symbol" xorm:"symbol"`
	SourceID      int             `json:"source_id" xorm:"source_id"`
	EnableTrade   SymbTradeRight  `json:"enable_trade" xorm:"enable_trade"`
	Leverage      int32           `json:"leverage" xorm:"-"`
	SecurityID    int             `json:"security_id" xorm:"security_id"`
	MarginInitial decimal.Decimal `json:"margin_initial" xorm:"margin_initial"`
	MarginDivider decimal.Decimal `json:"margin_divider" xorm:"margin_divider"`
	Percentage    decimal.Decimal `json:"percentage" xorm:"percentage"`
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

type SymbTradeRight int

const (
	NotSupport SymbTradeRight = iota
	OnlyClose
	OpenClose
)

type HolidayCategory int

const (
	HolidayAll HolidayCategory = iota
	HolidaySecurity
	HolidaySymbol
	HolidaySource
)

type TimeDay int

const (
	DayDIff TimeDay = iota
	DaySameTimeSame
	DaySameTimeDIff
)
