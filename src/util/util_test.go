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
	strs := table.Sprintf()

	if len(strs) != 3 {
		t.Log(strs)
		t.Error("util.Table.Sprintf() 生成的条目数量错误")
	}
	if strs[0] != "name    age   " {
		t.Log(strs[0])
		t.Error("util.Table.Sprintf() 生成的Lable有误")
	}

	if strs[1] != "John    30    " {
		t.Log(strs[0])
		t.Error("util.Table.Sprintf() 生成的Lable有误")
	}

	if strs[2] != "Janeh   25    " {
		t.Log(strs[0])
		t.Error("util.Table.Sprintf() 生成的Lable有误")
	}
}
