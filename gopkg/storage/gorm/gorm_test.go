package gorm

import (
	"gorm.io/gorm/clause"
	"testing"
)

const (
	source = "root:root0987@tcp(127.0.0.1:3306)/test-db?parseTime=True"
)

func Test1(t *testing.T) {
	db, err := NewDb(NewDail1(source))
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&CreditCard{})
	if err != nil {
		t.Fatal(err)
	}

	var cards []CreditCard

	result := db.Debug().Model(&cards).Clauses(clause.Returning{}).Where("user_id = ?", 1).Update("number", 1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("OK")
}

func Test2(t *testing.T) {
	db, err := NewDb(NewDail1(source))
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&CreditCard{})
	if err != nil {
		t.Fatal(err)
	}

	var card CreditCard

	result := db.Debug().Model(card).Where("user_id = ?", 2).Update("number", 1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	t.Log(result.RowsAffected)
}
