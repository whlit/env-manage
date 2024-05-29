package util

import (
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	Columns  []string
	Selected func(map[string]string) bool
	data     []map[string]string
	lens     map[string]int
}

// 新建表格
func NewTable(columns []string) *Table {
	return &Table{
		Columns: columns,
		data:    make([]map[string]string, 0),
		lens:    make(map[string]int),
	}
}

// 打印
func (t *Table) Printf() {
	strs := t.Sprintf()
	for _, str := range strs {
		fmt.Println(str)
	}
}

// 新增数据
func (t *Table) Add(rows ...map[string]string) *Table {
	if t.lens == nil {
		t.lens = make(map[string]int)
	}
	for _, row := range rows {
		for _, column := range t.Columns {
			t.lens[column] = max(t.lens[column], len(row[column]))
		}
	}
	t.data = append(t.data, rows...)
	return t
}

// 获取格式化数据
func (t *Table) Sprintf() []string {
	var res []string = make([]string, len(t.data)+1)

	// 打印列名
	if t.Columns == nil {
		return res
	}
	var builder strings.Builder
	if t.Selected != nil {
		builder.WriteString("   ")
	}
	var formats map[string]string = make(map[string]string)
	for _, column := range t.Columns {
		formats[column] = strings.Join([]string{"%-", strconv.Itoa(max(t.lens[column], len(column))), "s   "}, "")
		builder.WriteString(fmt.Sprintf(formats[column], column))
	}
	res[0] = builder.String()
	// 打印数据
	if t.data == nil {
		return res
	}
	for i, row := range t.data {
		builder.Reset()
		// 是否打印选中标记
		if t.Selected != nil {
			if t.Selected(row) {
				builder.WriteString(" * ")
			} else {
				builder.WriteString("   ")
			}
		}
		for _, column := range t.Columns {
			builder.WriteString(fmt.Sprintf(formats[column], row[column]))
		}
		res[i+1] = builder.String()
	}
	return res
}
