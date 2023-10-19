package struct2sql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
)

// 定义接口 parser，用于解析结构体字段
type parser interface {
	parse(key reflect.StructField, val reflect.Value) (err error)
}

// 解析函数，接收一个实现 parser 接口的解析器和结构体 s，返回解析错误
func parse(p parser, s interface{}) error {
	keys := reflect.TypeOf(s)
	values := reflect.ValueOf(s)
	return recursion(p, keys, values)
}

// 递归函数，接收一个解析器 p，结构体类型 keys，结构体值 values，返回解析错误
func recursion(p parser, keys reflect.Type, values reflect.Value) error {
	keys = typeElem(keys)
	values = valueElem(values)
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Field(i)
		value := values.Field(i)
		if key.Type.Kind() == reflect.Struct {
			if err := recursion(p, key.Type, value); err != nil {
				return err
			}
		} else {
			err := p.parse(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 处理类型的元素，用于处理指针类型
func typeElem(key reflect.Type) reflect.Type {
	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}
	return key
}

// 处理值的元素，用于处理指针类型
func valueElem(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func parseTable(key reflect.StructField, val reflect.Value) (table string, exist bool, err error) {
	table, exist = key.Tag.Lookup("table")
	if !exist {
		return
	}
	if table == "" {
		err = fmt.Errorf("table tag value can not be empty")
		return
	}
	return
}

func parseInsertColumnAndValue(key reflect.StructField, val reflect.Value) (column string, exist bool, value interface{}, err error) {
	column, exist = key.Tag.Lookup("insert")
	if !exist {
		return
	}
	if column == "" {
		err = fmt.Errorf("column tag value can not be empty")
		return
	}
	column, omitempty := parseInsertColumn(column)
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
	value = val.Interface()
	return
}

func parseInsertColumn(tag string) (column string, omitempty bool) {
	s2 := strings.SplitN(tag, ",", 2)
	column = s2[0]
	if len(s2) == 2 && s2[1] == "omitempty" {
		omitempty = true
	}
	return
}

func parseSelectColumn(key reflect.StructField, val reflect.Value) (column string, exist bool, err error) {
	column, exist = key.Tag.Lookup("select")
	if !exist {
		return
	}
	if column == "" {
		err = fmt.Errorf("column tag value can not be empty")
		return
	}
	return
}

func parseUpdateColumnAndValue(key reflect.StructField, val reflect.Value) (column string, exist bool, value interface{}, err error) {
	tag, exist := key.Tag.Lookup("update")
	if !exist {
		return
	}
	// 如果当前值不能转换成interface，返回错误信息
	if !val.CanInterface() {
		err = fmt.Errorf("condition field %s can not interfaced", key.Name)
		return
	}
	column, omitempty := parseUpdateColumn(tag)
	if column == "" {
		err = fmt.Errorf("column tag value can not be empty")
		return
	}
	// 如果忽略空，并且val是对应的零值 直接返回
	if omitempty && val.IsZero() {
		return
	}
	value = val.Interface()
	return
}

func parseUpdateColumn(tag string) (column string, omitempty bool) {
	s2 := strings.SplitN(tag, ",", 2)
	column = s2[0]
	if len(s2) == 2 && s2[1] == "omitempty" {
		omitempty = true
	}
	return
}

// parse方法用于解析结构体中的where标签，并将解析结果添加到条件表达式中。
func parseCondition(key reflect.StructField, val reflect.Value) (condition squirrel.Sqlizer, exist bool, err error) {
	tag, exist := key.Tag.Lookup("where")
	if !exist {
		return
	}
	// 如果当前值不能转换成interface，返回错误信息
	if !val.CanInterface() {
		err = fmt.Errorf("condition field %s can not interfaced", key.Name)
		return
	}
	column, expression, omitempty := parseConditionColumn(tag)
	// 如果忽略空，并且val是对应的零值 直接返回
	if omitempty && val.IsZero() {
		return
	}
	switch {
	case strings.Contains(expression, "like"):
		condition = squirrel.Expr(column+" "+expression+" ?", fmt.Sprintf("%%%s%%", val.Interface()))
	case strings.Contains(expression, "plike"):
		condition = squirrel.Expr(column+" "+expression+" ?", fmt.Sprintf("%%%s", val.Interface()))
	case strings.Contains(expression, "slike"):
		condition = squirrel.Expr(column+" "+expression+" ?", fmt.Sprintf("%s%%", val.Interface()))
	default:
		condition = squirrel.Expr(column+" "+expression+" ?", val.Interface())
	}
	return
}

func parseConditionColumn(tag string) (column string, expression string, omitempty bool) {
	s3 := strings.SplitN(tag, ",", 3)
	column = s3[0]
	if len(s3) == 3 {
		expression = s3[1]
		if s3[2] == "omitempty" {
			omitempty = true
		}
	} else if len(s3) == 2 {
		expression = s3[1]
	} else {
		expression = "="
	}
	return
}

func parseOrder(key reflect.StructField, val reflect.Value) (order string, exist bool, err error) {
	order, exist = key.Tag.Lookup("order")
	if !exist {
		return
	}
	if order == "" {
		err = fmt.Errorf("order tag value can not be empty")
		return
	}
	return
}
func parseOffset(key reflect.StructField, val reflect.Value) (offset uint64, exist bool, err error) {
	_, exist = key.Tag.Lookup("offset")
	if !exist {
		return
	}
	if !val.CanInterface() {
		err = fmt.Errorf("offset value can not convert to interface")
		return
	}
	value := val.Interface()
	offset, ok := value.(uint64)
	if !ok {
		err = fmt.Errorf("offset type is not expected, expect uint64, given %v", val.Type())
	}
	return
}

func parseLimit(key reflect.StructField, val reflect.Value) (limit uint64, exist bool, err error) {
	_, exist = key.Tag.Lookup("limit")
	if !exist {
		return
	}
	if !val.CanInterface() {
		err = fmt.Errorf("limit value can not convert to interface")
		return
	}
	value := val.Interface()
	limit, ok := value.(uint64)
	if !ok {
		err = fmt.Errorf("limit type is not expected, expect uint64, given %v", val.Type())
	}
	return
}
