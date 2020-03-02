package mt4SyncEt6

import (
	"encoding/json"
	"fmt"
	"mt4SyncEt6/holiday"
	"testing"
	"time"
)

func TestHoliday(t *testing.T) {
	//Xorm
	engine := GetEngine()
	//Prepare json
	holidaySlice := make([]HolidayMt4, 300)
	err := json.Unmarshal([]byte(holiday.JsonHoliday), &holidaySlice)
	if err != nil {
		panic(err)
	}
	//process
	for _, v := range holidaySlice {
		fmt.Println("============start")
		if v.ID != 0 {
			//Init
			holidayEt6 := Holiday{
				ID:          v.ID,
				Date:        v.Date,
				Enable:      true,
				Category:    HolidaySecurity,
				Symbol:      v.Symbol,
				Description: v.Description,
			}
			holidayEt6.Enable = true

			//from and to
			//Conv  GMT+3
			fStrGMT3 := TimeConv(v.Date, v.From)
			tStrGMT3 := TimeConv(v.Date, v.To)
			//From and to is day day
			fUTC := TimeConvToUTC(fStrGMT3)
			tUTC := TimeConvToUTC(tStrGMT3)
			fmt.Printf("fUTC:%+v\n", fUTC)
			fmt.Printf("fStrGMT3:%+v\n", fStrGMT3)

			fmt.Printf("tUTC:%+v\n", tUTC)
			fmt.Printf("tStrGMT3:%+v\n", tStrGMT3)
			//From and to are not same day
			flag, isSame := IsSameDay(fUTC, tUTC)
			fmt.Println("isSame", isSame, flag)

			if !isSame {
				if flag == DayDIff {
					ff, ft, datef := ConvFromItem(fUTC)
					holidayEt6.Date = datef
					holidayEt6.From = ff
					holidayEt6.To = ft
					fmt.Printf("ff, ft, datef :==== %+v  %+v  %+v", ff, ft, datef)

					fmt.Printf("holidayEt6====:%+v", holidayEt6)
					_, err := engine.Table("holiday").Insert(holidayEt6)
					if err != nil {
						panic(err)
					}
					tf, tt, datet := ConvToItem(tUTC)
					fmt.Printf("tf, tt, datet====:%+v :%+v :%+v", tf, tt, datet)

					holidayEt6.Date = datet
					holidayEt6.From = tf
					holidayEt6.To = tt
					fmt.Printf("holidayEt6====:%+v", holidayEt6)
					_, err = engine.Table("holiday").Insert(holidayEt6)
					if err != nil {
						panic(err)
					}
					fmt.Println("============End")
					continue
				}

				if flag == DaySameTimeSame {
					ff, ft, datef := ConvFromItem(fUTC)
					holidayEt6.Date = datef
					holidayEt6.From = ff
					holidayEt6.To = ft
					fmt.Printf("ff, ft, datef :==== %+v  %+v  %+v", ff, ft, datef)

					fmt.Printf("holidayEt6====:%+v", holidayEt6)
					_, err := engine.Table("holiday").Insert(holidayEt6)
					if err != nil {
						panic(err)
					}
					tf, tt, datet := ConvToItemAddOneDay(tUTC)
					fmt.Printf("tf, tt, datet====:%+v :%+v :%+v", tf, tt, datet)

					holidayEt6.Date = datet
					holidayEt6.From = tf
					holidayEt6.To = tt
					fmt.Printf("holidayEt6====:%+v", holidayEt6)
					_, err = engine.Table("holiday").Insert(holidayEt6)
					if err != nil {
						panic(err)
					}
					fmt.Println("============End")
					continue
				}

			}
			fStr, _ := ConvertInsameSay(fUTC)
			tStr, tdate := ConvertInsameSay(tUTC)
			holidayEt6.Date = tdate
			holidayEt6.From = fStr
			holidayEt6.To = tStr
			fmt.Printf("holidayEt6====:%+v", holidayEt6)
			_, err := engine.Table("holiday").Insert(holidayEt6)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("============End")
	}
}

func TestTimeConvFrom(t *testing.T) {
	t1 := TimeConvToUTC("2019-12-31 14:04:05")
	t2 := TimeConvToUTC("2006-02-02 2:04:05")
	tStrGMT3 := TimeConv("2020-01-20", 1260)
	fmt.Println(IsSameDay(t1, t2))
	fmt.Println(tStrGMT3)

	d, _ := time.ParseDuration("24h")
	d1 := t1.Add(d)
	fmt.Println(t1.Date())
	fmt.Println(d1.Date())
}
