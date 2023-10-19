package struct2sql

import (
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
)

type insertconverter struct {
	table   string
	columns []string
	values  []interface{}
}

func (c *insertconverter) parse(key reflect.StructField, val reflect.Value) (err error) {
	// table
	if table, exist, serr := parseTable(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.table = table
	}
	// column and value
	if column, exist, value, serr := parseInsertColumnAndValue(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.columns = append(c.columns, column)
		c.values = append(c.values, value)
	}
	return
}

func (c *insertconverter) toSql() (sql string, data []interface{}, err error) {
	// table check
	if c.table == "" {
		err = fmt.Errorf("table can not be empty")
		return
	}

	// column length check
	if len(c.columns) == 0 {
		err = fmt.Errorf("no insert columns")
		return
	}
	// generate sql
	sql, data, err = squirrel.Insert(c.table).Columns(c.columns...).Values(c.values...).ToSql()
	return
}
