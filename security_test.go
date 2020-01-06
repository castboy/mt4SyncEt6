package mt4SyncEt6

import (
	"mt4SyncEt6/decimal"
	"mt4SyncEt6/security"
	"strconv"
	"testing"
	"fmt"
)

func TestGetGroup(t *testing.T) {
	for _, v := range security.Groups {
		k := " " + v + " "
		grp := security.GetStringMap(k)

		et6Group := AccountGroup{}
		et6Group.ID, _ = strconv.Atoi(grp["id"])
		et6Group.Name = v
		et6Group.DepositCurrency = grp["deposit_currency"]
		et6Group.MarginStopOut, _ = decimal.NewFromString(grp["margin_stop_out"])
		mm,_ := strconv.Atoi(grp["hedge_largeleg"])
		et6Group.MarginMode=MarginCalcMode(mm)

		GroupToDB(et6Group)

		t.Log(v)
	}
}

func TestGetSecurity(t *testing.T) {
	engine:= GetEngine()
	for _, v := range security.Groups {
		k := " " + v + " "
		secs := security.GetStringMap(k)

		for k, _ := range secs {
			if k == "id" || k == "deposit_currency" || k == "margin_stop_out" || k == "hedge_largeleg" {
				continue
			}

			sec := security.GetStringMap(" manager " + "." + k)

			et6Sec := ConGroupSec{}
			et6Sec.GroupId,_ = strconv.Atoi(secs["id"])
			et6Sec.EnableSecurity, _ = strconv.Atoi(sec["enable"])
			et6Sec.EnableTrade, _ = strconv.Atoi(sec["trade"])
			et6Sec.LotMin, _ = decimal.NewFromString(sec["lot_min"])
			et6Sec.LotMax, _ = decimal.NewFromString(sec["lot_max"])
			et6Sec.LotStep, _ = decimal.NewFromString(sec["lot_step"])
			et6Sec.SpreadDiff, _ = strconv.Atoi(sec["spread_diff"])
			et6Sec.Commission, _ = decimal.NewFromString(sec["comm_base"])

			sql := fmt.Sprintf("select `id` FROM security WHERE `security_name`='%s'", k)
			row ,_:= engine.QueryString(sql)

			et6Sec.SecurityId,_ = strconv.Atoi(row[0]["id"])

			_, err := engine.Table("con_group_sec").Insert(et6Sec)
			if err != nil {
				fmt.Println(err)
			}
			t.Log(k, et6Sec)
		}
	}
}

