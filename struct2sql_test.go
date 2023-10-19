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

type TestSelectStruct struct {
	ID       int    ` select:"id" where:"name,=,omitempty"`
	Name     string ` select:"name" where:"name,like,omitempty"`
	Age      int    ` select:"age"`
	IsActive bool   ` select:"is_active"`
	_        string `table:"my_table"`
}

func TestConvert2SelectSql(t *testing.T) {
	testStruct := TestSelectStruct{}
	sql, _, err := Convert2SelectSql(testStruct)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证生成的 SQL 语句
	expectedSQL := "SELECT id,name,age,is_active FROM my_table"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, sql)
	}
}

func TestConvert2SelectSqlWithConditions(t *testing.T) {
	testStruct := TestSelectStruct{
		ID:       1,
		Name:     "John",
		Age:      30,
		IsActive: true,
	}
	sql, data, err := Convert2SelectSql(testStruct)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证生成的 SQL 语句
	expectedSQL := "SELECT id,name,age,is_active FROM my_table WHERE name = ? AND name like ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, sql)
	}

	// 验证生成的数据
	expectedData := []interface{}{1, "%John%"}
	for i := 0; i < len(expectedData); i++ {
		if data[i] != expectedData[i] {
			t.Errorf("Expected data: %v, but got: %v", expectedData, data)
			break
		}
	}
}

func TestConvert2SelectSqlWithOrderAndLimitOffset(t *testing.T) {
	testStruct := struct {
		TestSelectStruct
		_      string `order:"id desc, name"`
		Offset uint64 `offset:""`
		Limit  uint64 `limit:""`
	}{
		TestSelectStruct: TestSelectStruct{},
		Offset:           10,
		Limit:            20,
	}
	sql, _, err := Convert2SelectSql(testStruct)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证生成的 SQL 语句
	expectedSQL := "SELECT id,name,age,is_active FROM my_table ORDER BY id desc, name LIMIT 20 OFFSET 10"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, sql)
	}
}
