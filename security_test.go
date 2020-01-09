package mt4SyncEt6

import (
	"fmt"
	"mt4SyncEt6/decimal"
	"mt4SyncEt6/security"
	"strconv"
	"testing"
)

func TestGetGroup(t *testing.T) {
	engine := GetEngine()
	for _, v := range security.Groups {
		k := " " + v + " "
		grp := security.GetStringMap(k)

		et6Group := AccountGroup{}
		et6Group.ID, _ = strconv.Atoi(grp["id"])
		et6Group.Name = v
		et6Group.DepositCurrency = grp["deposit_currency"]
		et6Group.MarginStopOut, _ = decimal.NewFromString(grp["margin_stop_out"])
		mm, _ := strconv.Atoi(grp["hedge_largeleg"])
		et6Group.MarginMode = MarginCalcMode(mm)

		_, err := engine.Table("account_group").Insert(et6Group)
		if err != nil {
			fmt.Println(err)
		}

		t.Log(v,et6Group.ID)
	}
}

func TestGetSecurity(t *testing.T) {
	engine := GetEngine()
	for _, v := range security.Groups {
		k := " " + v + " "
		secs := security.GetStringMap(k)

		for k, _ := range secs {
			if k == "id" || k == "deposit_currency" || k == "margin_stop_out" || k == "hedge_largeleg" {
				continue
			}

			sec := security.GetStringMap(" manager " + "." + k)

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

			sql := fmt.Sprintf("select `id` FROM security WHERE `security_name`='%s'", k)
			row, _ := engine.QueryString(sql)

			et6Sec.SecurityId, _ = strconv.Atoi(row[0]["id"])

			_, err := engine.Table("con_group_sec").Insert(et6Sec)
			if err != nil {
				fmt.Println(err)
			}
			t.Log(k, et6Sec)
		}
	}
}
