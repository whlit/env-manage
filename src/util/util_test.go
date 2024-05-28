package util

import (
	"testing"
)

func TestTablePrint(t *testing.T) {
	table := &Table{
		Columns: []string{"name", "age"},
	}
	table.Add(map[string]string{
		"name": "John",
		"age":  "30",
	})
	table.Add(map[string]string{
		"name": "Janeh",
		"age":  "25",
	})
	table.Print()
}

func TestGetExeName(t *testing.T) {
	name := GetExeName()
	t.Log(name)
}
