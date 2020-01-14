package mt4SyncEt6

import (
	"encoding/json"
	"fmt"
	currency "mt4SyncEt6/currency"
	"testing"
)

func TestCurrency(t *testing.T) {
	var cfhQuoteSymbolMap = map[string]string{
		// transfer from CFH to TW
		"UKOUSD": "XBRUSD",
		"USOUSD": "XTIUSD",
		"E50EUR": "EUSTX50",
		"D30EUR": "GER30",
		"H33HKD": "HK50",
		"225JPY": "JPN225",
		"NASUSD": "NAS100",
		"U30USD": "US30",
		"SPXUSD": "USA500",
		"E35EUR": "SPA35",
		"100GBP": "UK100",
		"200AUD": "AUS200",
		"F40EUR": "FRA40",
	}

	str, _ := currency.GetCurrency()
	var currencies []CurrencyInfo
	json.Unmarshal([]byte(str), &currencies)

	engine := GetEngine()
	for _, v := range currencies {
		if s, ok := cfhQuoteSymbolMap[v.Symbol]; ok {
			v.Symbol = s
		}
		sql := "update source set currency = ? where source = ?"
		_, err := engine.Exec(sql, v.Currency,  v.Symbol)
		if err != nil {
			fmt.Println(err)
		}
	}
}
