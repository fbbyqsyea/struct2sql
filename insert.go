package struct2sql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
)

type insertconverter struct {
	table   string
	columns []string
	values  []interface{}
}

func (c *insertconverter) parse(key reflect.StructField, val reflect.Value) (err error) {
	if tag := key.Tag.Get("table"); tag != "" {
		if c.table != "" {
			err = fmt.Errorf("table tag duplicate definition")
			return
		}
		c.table = tag
	}
	if tag := key.Tag.Get("insert"); tag != "" {
		column, omitempty := c.parseInsertColumn(tag)
		if column == "" {
			err = fmt.Errorf("field %s 's insert column can not be empty", key.Name)
			return
		}
		if !val.CanInterface() {
			err = fmt.Errorf("field %s 's value can not interfaced", key.Name)
			return
		}
		if !omitempty && val.IsZero() {
			err = fmt.Errorf("field %s 's insert value can not be empty", key.Name)
			return
		}
		c.columns = append(c.columns, column)
		c.values = append(c.values, val.Interface())
	}
	return
}

func (c *insertconverter) parseInsertColumn(it string) (column string, omitempty bool) {
	s2 := strings.SplitN(it, ",", 2)
	column = s2[0]
	if len(s2) == 2 && s2[1] == "omitempty" {
		omitempty = true
	}
	return
}
func (c *insertconverter) toSql() (sql string, data []interface{}, err error) {

	if c.table == "" {
		err = fmt.Errorf("table can not be empty")
		return
	}

	if len(c.columns) == 0 {
		err = fmt.Errorf("no insert columns")
		return
	}

	sql, data, err = squirrel.Insert(c.table).Columns(c.columns...).Values(c.values...).ToSql()
	return
}
