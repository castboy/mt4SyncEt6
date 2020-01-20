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

	//oldId := make(map[string]int)
	//newId := make(map[string]int)
	//sql := "select `id` ,`name` from account_group"
	//row, _ := engine.QueryString(sql)
	//for _, v := range row {
	//	k := v[`name`]
	//	oldId[" "+k+" "], _ = strconv.Atoi(v[`id`])
	//}

	sql := "Truncate Table account_group"
	_, err := engine.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}

	for _, group := range account_group.Groups {
		grp := con_group_sec.GetSecStringMap(group)
		if len(grp) == 0 {
			group = " " + group + " "
			grp = con_group_sec.GetSecStringMap(group)
		}

		et6Group := AccountGroup{}
		et6Group.ID, _ = strconv.Atoi(grp["id"])
		et6Group.Name = group
		et6Group.DepositCurrency = grp["deposit_currency"]
		et6Group.MarginStopOut, _ = decimal.NewFromString(grp["margin_stop_out"])
		mm, _ := strconv.Atoi(grp["hedge_largeleg"])
		et6Group.MarginMode = MarginCalcMode(mm)

		_, err := engine.Table("account_group").Insert(et6Group)
		if err != nil {
			fmt.Println(err)
		}

		//newId[k] = et6Group.ID
		t.Log(group, et6Group.ID)
	}

	//for k, v := range oldId {
	//	if newId[k] != v {
	//		sql = fmt.Sprintf("UPDATE account SET group_id = %d WHERE group_id = %d", newId[k], oldId[k])
	//		_, err := engine.Exec(sql)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	//}
}
