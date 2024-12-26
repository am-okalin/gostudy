package gorm

import (
	"database/sql"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

// User 有多张 CreditCard，UserID 是外键
type User struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func NewDail1(dataSource string) gorm.Dialector {
	return mysql.Open(dataSource)
}

func NewDail2(dataSource string) (gorm.Dialector, error) {
	sqlDB, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	dial := mysql.New(mysql.Config{
		Conn: sqlDB,
		//DSN: dataSource, //直接配置这个就不用配置conn了
	})
	return dial, nil
}

func NewDb(dial gorm.Dialector) (*gorm.DB, error) {
	return gorm.Open(dial, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
