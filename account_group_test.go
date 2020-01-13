package mt4SyncEt6

import (
	"fmt"
	"mt4SyncEt6/account_group"
	"mt4SyncEt6/con_group_sec"
	"mt4SyncEt6/decimal"
	"strconv"
	"testing"
)

func TestAccountGroup(t *testing.T) {
	engine := GetEngine()

	sql := "Truncate Table account_group"
	_, err := engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range account_group.Groups {
		k := " " + v + " "
		grp := con_group_sec.GetSecStringMap(k)

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

		t.Log(v, et6Group.ID)
	}
}
