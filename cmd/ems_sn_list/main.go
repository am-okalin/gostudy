package main

import (
	"encoding/csv"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type EmsEquipmentValue struct {
	ID          int    `gorm:"primarykey"`
	EquipmentId int    `gorm:"not null;comment:设备ID;Index"`
	Attribute   string `gorm:"not null;size:255;comment:设备属性[sn, imei1, imei2, fsn];Index"`
	Value       string `gorm:"not null;size:255;comment:设备属性值;uniqueIndex"`
}

func main() {
	// db配置
	dsn := "root:root0987@tcp(127.0.0.1:3306)/test-db?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 输出文件
	ofile, err := os.Create("./ems_sn_list/ems_fsn.csv")
	if err != nil {
		panic(err)
	}
	defer ofile.Close()
	fmt.Fprintf(ofile, "id,sn,fsn\n")

	//读取文件
	file, err := os.Open("./ems_sn_list/sns.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建带缓冲的读取器
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	sns := make([]string, 0, len(lines))
	for i := range lines {
		if lines[i][0] != "" {
			sns = append(sns, lines[i][0])
		}
	}

	for i := 0; i < len(sns); i += 1000 {
		var end int
		if i+1000 < len(sns) {
			end = i + 1000
		} else {
			end = len(sns)
		}
		ListFsn(db, ofile, sns[i:end])
	}
}

type Obj struct {
	ID  int
	Sn  string
	Fsn string
}

func ListFsn(db *gorm.DB, ofile *os.File, sns []string) {
	objs := make([]Obj, 0, len(sns))
	sql := `select e.id, sn, ev.value as fsn from ems_equipments e left join ems_equipment_values ev on e.id = ev.equipment_id and ev.attribute="fsn" where sn in (?)`
	db.Raw(sql, sns).Scan(&objs)

	//将结果集输出到 esm_fsn.csv
	for _, obj := range objs {
		fmt.Fprintf(ofile, "%d,%s,%s\n", obj.ID, obj.Sn, obj.Fsn)
	}
}
