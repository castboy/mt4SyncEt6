package source_compare_dev_pro

import (
	"mt4SyncEt6"
)

func SourceCompare() bool {
	//Data prepare
	sourceItemDev := new([]mt4SyncEt6.Source)
	sourceItemPro := new([]mt4SyncEt6.Source)
	//DB driver prepare
	devXrom, _ := mt4SyncEt6.NewET6EngineXorm()
	proXrom, _ := mt4SyncEt6.NewProduceEngineXorm()

	//GetData
	devXrom.Table("source").Find(&sourceItemDev)
	proXrom.Table("source").Find(&sourceItemPro)

	return true
}
