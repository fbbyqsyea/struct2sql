package struct2sql

import (
	"reflect"
)

type selector struct {
	columns []string
	values  []interface{}
}

func (c selector) parse(key reflect.StructField, val reflect.Value) (sql string, data []interface{}, err error) {

	return
}
