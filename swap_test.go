package mt4SyncEt6

import (
	"encoding/json"
	"fmt"
	"mt4SyncEt6/swap"
	"testing"
)

func TestSwap(t *testing.T) {
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

	str, _ := swap.GetSwap()
	var swaps []SwapInfo
	json.Unmarshal([]byte(str), &swaps)

	engine := GetEngine()
	for _, v := range swaps {
		if s, ok := cfhQuoteSymbolMap[v.Symbol]; ok {
			v.Symbol = s
		}
		sql := "update source set swap_long = ?,swap_short = ?,source_cn = ?,swap_3_day = ? where source = ?"
		_, err := engine.Exec(sql, v.SwapLong, v.SwapShort, v.SourceCN,v.Swap3Day, v.Symbol)
		if err != nil {
			fmt.Println(err)
		}
	}
}
