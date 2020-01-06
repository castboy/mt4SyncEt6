package mt4SyncEt6

import (
	"mt4SyncEt6/decimal"
	"mt4SyncEt6/security"
	"strconv"
	"testing"
)

func TestGetGroup(t *testing.T) {
	for _, v := range security.Groups {
		k := " " + v + " "
		grp := security.GetGroup(k)

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
	et6 := security.GetSecurity()
	t.Log(et6)
}

