package struct2sql

import (
	"reflect"
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

func parseTable(key reflect.StructField, val reflect.Value) (table string) {
	table = key.Tag.Get("table")
	return
}
