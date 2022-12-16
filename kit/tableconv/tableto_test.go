package tableconv

import (
	"testing"
)

type User struct {
	Id, Name, Age, CreateAt string
}

func MockTable() [][]string {
	//cl, rl := 4, 10
	return [][]string{
		{"Id", "Name", "Age", "CreateAt"},
		{"1", "a", "10", "2010-01-01 12:00:00"},
		{"2", "b", "20", "2000-01-01 12:00:00"},
		{"3", "c", "30", "1990-01-01 12:00:00"},
		{"4", "d", "40", "1980-01-01 12:00:00"},
	}
}

func TestTo(t *testing.T) {
	table := MockTable()
	var objs []User
	err := ToObj(table, &objs)
	if err != nil {
		t.Error(err)
	}

	t.Log("done")
}
