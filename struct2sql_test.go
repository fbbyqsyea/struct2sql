package struct2sql

import (
	"testing"
)

type TestStruct struct {
	ID   int    `insert:"id"`
	Name string `insert:"name"`
	Age  int    `insert:"age,omitempty"`
	_    string `table:"my_table"`
}

func TestConvert2InsertSql(t *testing.T) {
	testStruct := TestStruct{
		ID:   1,
		Name: "John",
		Age:  30,
	}

	sql, data, err := Convert2InsertSql(testStruct)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证生成的 SQL 语句
	expectedSQL := "INSERT INTO my_table (id,name,age) VALUES (?,?,?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, sql)
	}

	// 验证生成的数据
	expectedData := []interface{}{1, "John", 30}
	for i := 0; i < len(expectedData); i++ {
		if data[i] != expectedData[i] {
			t.Errorf("Expected data: %v, but got: %v", expectedData, data)
			break
		}
	}
}
