package struct2sql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
)

type updateconverter struct {
	table          string
	columnValueMap map[string]interface{}
	conditions     []squirrel.Sqlizer
	order          []string
	limit          uint64
}

func (c *updateconverter) parse(key reflect.StructField, val reflect.Value) (err error) {
	// table
	if table, exist, serr := parseTable(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.table = table
	}
	// cloumns
	if column, exist, value, serr := parseUpdateColumnAndValue(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.columnValueMap[column] = value
	}
	// conditions
	if condition, exist, serr := parseCondition(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.conditions = append(c.conditions, condition)
	}
	// order
	if order, exist, serr := parseOrder(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.order = append(c.order, order)
	}
	// limit
	if limit, exist, serr := parseLimit(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.limit = limit
	}
	return
}

func (c *updateconverter) toSql() (sql string, data []interface{}, err error) {
	// table check
	if c.table == "" {
		err = fmt.Errorf("table can not be empty")
		return
	}

	// column length check
	if len(c.columnValueMap) == 0 {
		err = fmt.Errorf("no update column")
		return
	}

	sb := squirrel.Update(c.table)

	// column value maps
	for column, value := range c.columnValueMap {
		sb = sb.Set(column, value)
	}

	// condition
	for _, condition := range c.conditions {
		sb = sb.Where(condition)
	}

	// limit
	if c.limit > 0 {
		sb = sb.Limit(c.limit)
	}

	// order by
	if len(c.order) > 0 {
		sb = sb.OrderBy(strings.Join(c.order, ","))
	}

	return sb.ToSql()
}
