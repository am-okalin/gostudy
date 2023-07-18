package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type EmsEquipment struct {
	ID int `gorm:"primarykey"`
}

type EmsEquipmentValue struct {
	ID          int    `gorm:"primarykey"`
	EquipmentId int    `gorm:"not null;comment:设备ID;Index"`
	Attribute   string `gorm:"not null;size:255;comment:设备属性[sn, imei1, imei2, fsn];Index"`
	Value       string `gorm:"not null;size:255;comment:设备属性值;uniqueIndex"`
}

func main() {
	//dsn := "root:root0987@tcp(127.0.0.1:3306)/test-db?parseTime=True"
	dsn := "root:ng6brE9LqSrt0M4H8HUD@tcp(ec2-18-167-6-45.ap-east-1.compute.amazonaws.com:3306)/fstln_apps?parseTime=True"
	//dsn := "rds_ems_admin:cVHRoeCoDjA9CXodtJCL@tcp(database-1.cjqgip5trpn2.eu-central-1.rds.amazonaws.com:3306)/fstln_ems?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//err = BeachCreateFsn(db, "./test.xlsx", "test.txt")
	//if err != nil {
	//	panic(err)
	//}

	err = BeachCreateFsn(db, "./NT-IN (1).xlsx", "in.txt")
	if err != nil {
		panic(err)
	}

	err = BeachCreateFsn(db, "./NT-XA.xlsx", "xa.txt")
	if err != nil {
		panic(err)
	}
}

func BeachCreateFsn(db *gorm.DB, input, output string) error {
	ifile, err := xlsx.OpenFile(input)
	if err != nil {
		panic(err)
	}

	ofile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer ofile.Close()

	for i, row := range ifile.Sheets[0].Rows {
		if i == 0 {
			continue
		}
		err = CreateFsn(db, row.Cells[1].Value, row.Cells[0].Value)
		if err != nil {
			fmt.Fprintf(ofile, "%d,%s,%s,%s\n", i, row.Cells[0].Value, row.Cells[1].Value, err.Error())
		}
	}
	return nil
}

func CreateFsn(db *gorm.DB, sn, fsn string) error {
	// 查询设备表, 获取equipment_id, 不存在则报错

	equipment := EmsEquipment{}
	sql := "select id from ems_equipments where sn=?"
	res := db.Raw(sql, sn).Scan(&equipment)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) || equipment.ID == 0 {
		return errors.New(fmt.Sprintf("sn:%s not found", sn))
	}
	// 查询设备eav表, fsn已存在则报错
	ev := EmsEquipmentValue{}
	sql = "select id, equipment_id, attribute, value from ems_equipment_values where equipment_id=? and attribute='fsn' or value=?"
	res = db.Raw(sql, equipment.ID, fsn).Scan(&ev)
	if ev.ID != 0 {
		return errors.New(fmt.Sprintf("fsn already existing, %v", ev))
	}
	// 执行更新操作
	ev = EmsEquipmentValue{
		EquipmentId: equipment.ID,
		Attribute:   "fsn",
		Value:       fsn,
	}
	res = db.Create(&ev)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
