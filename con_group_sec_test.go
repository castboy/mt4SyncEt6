package mt4SyncEt6

import (
	"fmt"
	"mt4SyncEt6/account_group"
	"mt4SyncEt6/con_group_sec"
	"mt4SyncEt6/decimal"
	"strconv"
	"testing"
)

func TestConGroupSec(t *testing.T) {
	engine := GetEngine()

	sql := "Truncate Table con_group_sec"
	_, err := engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}

	for _, group := range account_group.Groups {
		secs := con_group_sec.GetSecStringMap(group)
		if len(secs) == 0 {
			group = " " + group + " "
			secs = con_group_sec.GetSecStringMap(group)
		}

			for field, _ := range secs {
			if field == "id" || field == "deposit_currency" || field == "margin_stop_out" || field == "hedge_largeleg" {
				continue
			}

			sec := con_group_sec.GetSecStringMap(group + "." + field)

			et6Sec := ConGroupSec{}
			size := *decimal.NewDecFromInt(100)
			et6Sec.GroupId, _ = strconv.Atoi(secs["id"])
			et6Sec.EnableSecurity, _ = strconv.Atoi(sec["enable"])
			et6Sec.EnableTrade, _ = strconv.Atoi(sec["trade"])
			lotMin, _ := decimal.NewFromString(sec["lot_min"])
			lotMax, _ := decimal.NewFromString(sec["lot_max"])
			lotStep, _ := decimal.NewFromString(sec["lot_step"])
			et6Sec.LotMin = lotMin.Div(size)
			et6Sec.LotMax = lotMax.Div(size)
			et6Sec.LotStep = lotStep.Div(size)
			et6Sec.SpreadDiff, _ = strconv.Atoi(sec["spread_diff"])
			et6Sec.Commission, _ = decimal.NewFromString(sec["comm_base"])

			sql := fmt.Sprintf("select `id` FROM security WHERE `security_name`='%s'", field)
			row, err := engine.QueryString(sql)
			if row == nil {
				fmt.Println(field)
				continue
			}
			et6Sec.SecurityId, err = strconv.Atoi(row[0]["id"])

			_, err = engine.Table("con_group_sec").Insert(et6Sec)
			if err != nil {
				fmt.Println(err)
			}
			t.Log(field, et6Sec)
		}
	}
}
