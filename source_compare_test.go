package mt4SyncEt6

import (
	"fmt"
	"testing"
)

func Test_sourceCompare(t *testing.T) {
	//Prepare
	var dev []Source
	var pro []Source
	devXrom, err := NewET6EngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	proXrom, err := NewProduceEngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	err = devXrom.SQL("select * from source").Find(&dev)
	if err != nil {
		fmt.Println("err", err)
	}
	err = proXrom.SQL("select * from source").Find(&pro)

	if err != nil {
		fmt.Println("err", err)
	}
	for _, devOne := range dev {
		flag := false
		for _, proOne := range pro {
			if devOne.Source == proOne.Source {
				flag = true
				if devOne.SourceType != proOne.SourceType {
					fmt.Printf(" Not compare:sourceType  SourceName:%+v,  Not compare devOne.SourceType:%+v, proOne.SourceType:%+v\n", devOne.Source, devOne.SourceType, proOne.SourceType)
				}

				if devOne.Digits != proOne.Digits {
					fmt.Printf(" Not compare:Digits  SourceName:%+v,  Not compare devOne.Digits:%+v, proOne.Digits:%+v\n", devOne.Source, devOne.Digits, proOne.Digits)
				}

				if devOne.Multiply != proOne.Multiply {
					fmt.Printf(" Not compare:Multiply  SourceName:%+v,  Not compare devOne.Multiply:%+v, proOne.Multiply:%+v\n", devOne.Source, devOne.Multiply, proOne.Multiply)
				}

				if devOne.ContractSize != proOne.ContractSize {
					fmt.Printf(" Not compare:ContractSize  SourceName:%+v,  Not compare devOne.ContractSize:%+v, proOne.ContractSize:%+v\n", devOne.Source, devOne.ContractSize, proOne.ContractSize)
				}

				if devOne.StopsLevel != proOne.StopsLevel {
					fmt.Printf(" Not compare:StopsLevel  SourceName:%+v,  Not compare devOne.StopsLevel:%+v, proOne.StopsLevel:%+v\n", devOne.Source, devOne.StopsLevel, proOne.StopsLevel)
				}

				if devOne.ProfitMode != proOne.ProfitMode {
					fmt.Printf(" Not compare:ProfitMode  SourceName:%+v,  Not compare devOne.ProfitMode:%+v, proOne.ProfitMode:%+v\n", devOne.Source, devOne.ProfitMode, proOne.ProfitMode)
				}

				if devOne.ProfitCurrency != proOne.ProfitCurrency {
					fmt.Printf(" Not compare:ProfitCurrency  SourceName:%+v,  Not compare devOne.ProfitCurrency:%+v, proOne.ProfitCurrency:%+v\n", devOne.Source, devOne.ProfitCurrency, proOne.ProfitCurrency)
				}

				if devOne.MarginMode != proOne.MarginMode {
					fmt.Printf(" Not compare:MarginMode  SourceName:%+v,  Not compare devOne.MarginMode:%+v, proOne.MarginMode:%+v\n", devOne.Source, devOne.MarginMode, proOne.MarginMode)
				}

				if devOne.MarginCurrency != proOne.MarginCurrency {
					fmt.Printf(" Not compare:MarginCurrency  SourceName:%+v,  Not compare devOne.MarginCurrency:%+v, proOne.MarginCurrency:%+v\n", devOne.Source, devOne.MarginCurrency, proOne.MarginCurrency)
				}

				if devOne.SwapType != proOne.SwapType {
					fmt.Printf(" Not compare:SwapType  SourceName:%+v,  Not compare devOne.SwapType:%+v, proOne.SwapType:%+v\n", devOne.Source, devOne.SwapType, proOne.SwapType)
				}

				if devOne.SwapLong != proOne.SwapLong {
					fmt.Printf(" Not compare:SwapLong  SourceName:%+v,  Not compare devOne.SwapLong:%+v, proOne.SwapLong:%+v\n", devOne.Source, devOne.SwapLong, proOne.SwapLong)
				}

				if devOne.SwapShort != proOne.SwapShort {
					fmt.Printf(" Not compare:SwapShort  SourceName:%+v,  Not compare devOne.SwapShort:%+v, proOne.SwapShort:%+v\n", devOne.Source, devOne.SwapShort, proOne.SwapShort)
				}

				if devOne.SwapCurrency != proOne.SwapCurrency {
					fmt.Printf(" Not compare:SwapCurrency  SourceName:%+v,  Not compare devOne.SwapCurrency:%+v, proOne.SwapCurrency:%+v\n", devOne.Source, devOne.SwapCurrency, proOne.SwapCurrency)
				}

				if devOne.Swap3Day != proOne.Swap3Day {
					fmt.Printf(" Not compare:Swap3Day  SourceName:%+v,  Not compare devOne.Swap3Day:%+v, proOne.Swap3Day:%+v\n", devOne.Source, devOne.Swap3Day, proOne.Swap3Day)
				}
			}
			flag = true
		}
		if flag == false {
			fmt.Println("Can not find the source", devOne.Source)
		}
	}
}

func Test_accountGroupCompare(t *testing.T) {
	//Prepare
	var dev []AccountGroup
	var pro []AccountGroup
	devXrom, err := NewET6EngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	proXrom, err := NewProduceEngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	err = devXrom.SQL("select * from account_group").Find(&dev)
	if err != nil {
		fmt.Println("err", err)
	}
	err = proXrom.SQL("select * from account_group").Find(&pro)

	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("dev", dev)
	fmt.Println("pro", pro)
	for _, devOne := range dev {
		flag := false
		for _, proOne := range pro {
			if devOne.Name == proOne.Name {
				flag = true
				if devOne.Name != proOne.Name {
					fmt.Printf(" Not compare:Name  SourceName:%+v,  Not compare devOne.Name:%+v, proOne.Name:%+v\n", devOne.Name, devOne.Name, proOne.Name)
				}

				if devOne.DepositCurrency != proOne.DepositCurrency {
					fmt.Printf(" Not compare:DepositCurrency  SourceName:%+v,  Not compare devOne.DepositCurrency:%+v, proOne.DepositCurrency:%+v\n", devOne.Name, devOne.DepositCurrency, proOne.DepositCurrency)
				}

				if devOne.MarginStopOut != proOne.MarginStopOut {
					fmt.Printf(" Not compare:MarginStopOut  SourceName:%+v,  Not compare devOne.MarginStopOut:%+v, proOne.MarginStopOut:%+v\n", devOne.Name, devOne.MarginStopOut, proOne.MarginStopOut)
				}

				if devOne.MarginMode != proOne.MarginMode {
					fmt.Printf(" Not compare:MarginMode  SourceName:%+v,  Not compare devOne.MarginMode:%+v, proOne.MarginMode:%+v\n", devOne.Name, devOne.MarginMode, proOne.MarginMode)
				}
			}
		}
		if flag == false {
			fmt.Println("Not find in pro ", devOne.Name)
		}
	}
}

func Test_CongroupSecurityCompare(t *testing.T) {
	//Prepare
	var dev []ConGroupSec
	var pro []ConGroupSec
	devXrom, err := NewET6EngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	proXrom, err := NewProduceEngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	err = devXrom.SQL("select * from con_group_sec").Find(&dev)
	if err != nil {
		fmt.Println("err", err)
	}
	err = proXrom.SQL("select * from con_group_sec").Find(&pro)

	if err != nil {
		fmt.Println("err", err)
	}
	for _, devOne := range dev {
		flag := false
		for _, proOne := range pro {
			if devOne.GroupId == proOne.GroupId && devOne.SecurityId == proOne.SecurityId {
				flag = true
				if devOne.EnableSecurity != proOne.EnableSecurity {
					fmt.Printf(" Not compare:Name  ID:%+v,  Not compare devOne.Name:%+v, proOne.Name:%+v\n", devOne.ID, devOne.EnableSecurity, proOne.EnableSecurity)
				}

				if devOne.EnableTrade != proOne.EnableTrade {
					fmt.Printf(" Not compare:EnableTrade  ID:%+v,  Not compare devOne.EnableTrade:%+v, proOne.EnableTrade:%+v\n", devOne.ID, devOne.EnableTrade, proOne.EnableTrade)
				}

				if devOne.LotMin != proOne.LotMin {
					fmt.Printf(" Not compare:LotMin  ID:%+v,  Not compare devOne.LotMin:%+v, proOne.LotMin:%+v\n", devOne.ID, devOne.LotMin, proOne.LotMin)
				}

				if devOne.LotMax != proOne.LotMax {
					fmt.Printf(" Not compare:LotMax  ID:%+v,  Not compare devOne.LotMax:%+v, proOne.LotMax:%+v\n", devOne.ID, devOne.LotMax, proOne.LotMax)
				}

				if devOne.LotStep != proOne.LotStep {
					fmt.Printf(" Not compare:LotStep  ID:%+v,  Not compare devOne.LotStep:%+v, proOne.LotStep:%+v\n", devOne.ID, devOne.LotStep, proOne.LotStep)
				}

				if devOne.SpreadDiff != proOne.SpreadDiff {
					fmt.Printf(" Not compare:SpreadDiff  ID:%+v,  Not compare devOne.SpreadDiff:%+v, proOne.SpreadDiff:%+v\n", devOne.ID, devOne.SpreadDiff, proOne.SpreadDiff)
				}

				if devOne.Commission != proOne.Commission {
					fmt.Printf(" Not compare:Commission  ID:%+v,  Not compare devOne.Commission:%+v, proOne.Commission:%+v\n", devOne.ID, devOne.Commission, proOne.Commission)
				}
			}
		}
		if flag == false {
			fmt.Println("Not find in pro ", devOne.ID)
		}
	}
}

func Test_SessionCompare(t *testing.T) {
	//Prepare
	var dev []Session
	var pro []Session
	devXrom, err := NewET6EngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	proXrom, err := NewProduceEngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	err = devXrom.SQL("select * from session").Find(&dev)
	if err != nil {
		fmt.Println("err", err)
	}
	err = proXrom.SQL("select * from session").Find(&pro)

	if err != nil {
		fmt.Println("err", err)
	}
	for _, devOne := range dev {
		flag := false
		timeSpans := make([]string, 0)
		//Search time span from pro
		for _, proOne := range pro {
			if devOne.SourceID == proOne.SourceID && devOne.Type == proOne.Type  {
				timeSpans = append(timeSpans, proOne.TimeSpan)
			}
		}

		for _, span := range timeSpans {
			if span == devOne.TimeSpan {
				flag = true
				goto EXIT
			}
		}
	EXIT:
		if flag == false {
			fmt.Println("Can not find the source devOne ID:", devOne.ID)
		}
	}

}

func Test_SymbolCompare(t *testing.T) {
	var dev []Symbol
	var pro []Symbol
	devXrom, err := NewET6EngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	proXrom, err := NewProduceEngineXorm()
	if err != nil {
		fmt.Println("err", err)
	}
	err = devXrom.SQL("select * from symbol").Find(&dev)
	if err != nil {
		fmt.Println("err", err)
	}

	err = proXrom.SQL("select * from symbol").Find(&pro)
	if err != nil {
		fmt.Println("err", err)
	}

	for _, devOne := range dev {
		flag := false
		for _, proOne := range pro {
			if devOne.Symbol == proOne.Symbol {
				flag = true

				if devOne.EnableTrade != proOne.EnableTrade {
					fmt.Printf(" Not compare:EnableTrade  Name:%+v,  Not compare devOne.EnableTrade:%+v, proOne.EnableTrade:%+v\n", devOne.Symbol, devOne.EnableTrade, proOne.EnableTrade)
				}

				if devOne.Leverage != proOne.Leverage {
					fmt.Printf(" Not compare:Leverage  Name:%+v,  Not compare devOne.Leverage:%+v, proOne.Leverage:%+v\n", devOne.Symbol, devOne.Leverage, proOne.Leverage)
				}

				if devOne.SecurityID != proOne.SecurityID {
					fmt.Printf(" Not compare:SecurityID  Name:%+v,  Not compare devOne.SecurityID:%+v, proOne.SecurityID:%+v\n", devOne.Symbol, devOne.SecurityID, proOne.SecurityID)
				}

				if devOne.MarginInitial != proOne.MarginInitial {
					fmt.Printf(" Not compare:MarginInitial  Name:%+v,  Not compare devOne.MarginInitial:%+v, proOne.MarginInitial:%+v\n", devOne.Symbol, devOne.MarginInitial, proOne.MarginInitial)
				}

				if devOne.MarginDivider != proOne.MarginDivider {
					fmt.Printf(" Not compare:MarginDivider  Name:%+v,  Not compare devOne.MarginDivider:%+v, proOne.MarginDivider:%+v\n", devOne.Symbol, devOne.MarginDivider, proOne.MarginDivider)
				}

				if devOne.Percentage != proOne.Percentage {
					fmt.Printf(" Not compare:Percentage  Name:%+v,  Not compare devOne.Percentage:%+v, proOne.Percentage:%+v\n", devOne.Symbol, devOne.Percentage, proOne.Percentage)
				}
			}
			flag = true
		}
		if flag == false {
			fmt.Println("Can not find the source", devOne.Symbol)
		}
	}

}
