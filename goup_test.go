package mt4SyncEt6

import (
	"fmt"
	"testing"
)

func Test_Group_changeName(t *testing.T) {
	tsEngine := GetEngine()
	ugs := []*AccountGroup{}
	tsEngine.Table("account_group").Find(&ugs)
	for i := range ugs {
		fmt.Printf("%+v\n", ugs[i])
		ugs[i].ID = 0
		ugs[i].Name = ugs[i].Name + "_EVO"
		tsEngine.Table(AccountGroup{}).Insert(ugs[i])
	}
}

func Test_conGroup(t *testing.T) {
	tsEngine := GetEngine()
	ugs := []*AccountGroup{}
	conGroupSecs := []*ConGroupSec{}
	//获取原来用户组
	tsEngine.Table("account_group_old").Find(&ugs)
	tsEngine.Table("con_group_sec").Find(&conGroupSecs)
	for i := range ugs {
		for j := range conGroupSecs {
			conGroupSecs[j].ID = 0
			if ugs[i].ID == conGroupSecs[j].GroupId {
				//进行分组转换
				if ugs[i].ID == 3 {
					conGroupSecs[j].GroupId = 120
				}
				if ugs[i].ID == 5 {
					conGroupSecs[j].GroupId = 121
				}
				if ugs[i].ID == 6 {
					conGroupSecs[j].GroupId = 122
				}
				if ugs[i].ID == 7 {
					conGroupSecs[j].GroupId = 123
				}
				if ugs[i].ID == 8 {
					conGroupSecs[j].GroupId = 124
				}
				if ugs[i].ID == 32 {
					conGroupSecs[j].GroupId = 125
				}
				if ugs[i].ID == 33 {
					conGroupSecs[j].GroupId = 126
				}
				if ugs[i].ID == 34 {
					conGroupSecs[j].GroupId = 127
				}
				if ugs[i].ID == 35 {
					conGroupSecs[j].GroupId = 128
				}
				if ugs[i].ID == 36 {
					conGroupSecs[j].GroupId = 129
				}
				if ugs[i].ID == 37 {
					conGroupSecs[j].GroupId = 130
				}
				if ugs[i].ID == 38 {
					conGroupSecs[j].GroupId = 131
				}
				if ugs[i].ID == 39 {
					conGroupSecs[j].GroupId = 132
				}
				if ugs[i].ID == 40 {
					conGroupSecs[j].GroupId = 133
				}
				if ugs[i].ID == 41 {
					conGroupSecs[j].GroupId = 134
				}
				if ugs[i].ID == 42 {
					conGroupSecs[j].GroupId = 135
				}
				if ugs[i].ID == 43 {
					conGroupSecs[j].GroupId = 136
				}
				if ugs[i].ID == 109 {
					conGroupSecs[j].GroupId = 137
				}
				if ugs[i].ID == 110 {
					conGroupSecs[j].GroupId = 138
				}
			}
		}
	}
	//插入
	for i := range conGroupSecs {
		fmt.Printf("%+v\n", conGroupSecs[i])
		tsEngine.Table("con_group_sec").Insert(conGroupSecs[i])
	}
}
