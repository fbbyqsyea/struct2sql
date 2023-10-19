package struct2sql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
)

type selectconverter struct {
	table      string
	columns    []string
	conditions []squirrel.Sqlizer
	order      []string
	offset     uint64
	limit      uint64
}

func (c *selectconverter) parse(key reflect.StructField, val reflect.Value) (err error) {
	// table
	if table, exist, serr := parseTable(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.table = table
	}
	// cloumns
	if column, exist, serr := parseSelectColumn(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.columns = append(c.columns, column)
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
	// offset
	if offset, exist, serr := parseOffset(key, val); exist {
		if serr != nil {
			err = serr
			return
		}
		c.offset = offset
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

func (c *selectconverter) toSql() (sql string, data []interface{}, err error) {
	// table check
	if c.table == "" {
		err = fmt.Errorf("table can not be empty")
		return
	}

	// column length check
	if len(c.columns) == 0 {
		err = fmt.Errorf("no select columns")
		return
	}

	sb := squirrel.Select(strings.Join(c.columns, ",")).From(c.table)

	// condition
	for _, condition := range c.conditions {
		sb = sb.Where(condition)
	}

	// offset
	if c.offset > 0 {
		sb = sb.Offset(c.offset)
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
