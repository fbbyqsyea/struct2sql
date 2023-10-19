package struct2sql

import "github.com/Masterminds/squirrel"

func Convert2SelectSql(s interface{}) (sql string, data []interface{}, err error) {
	selectconverter := selectconverter{
		columns:    make([]string, 0),
		conditions: make([]squirrel.Sqlizer, 0),
		order:      make([]string, 0),
	}
	err = parse(&selectconverter, s)
	if err != nil {
		return
	}
	sql, data, err = selectconverter.toSql()
	return
}

func Convert2InsertSql(s interface{}) (sql string, data []interface{}, err error) {
	insertconverter := insertconverter{
		columns: make([]string, 0),
		values:  make([]interface{}, 0),
	}
	err = parse(&insertconverter, s)
	if err != nil {
		return
	}
	sql, data, err = insertconverter.toSql()
	return
}

func Convert2UpdateSql(s interface{}) (sql string, data []interface{}, err error) {
	updateconverter := updateconverter{
		columnValueMap: make(map[string]interface{}),
		conditions:     make([]squirrel.Sqlizer, 0),
		order:          make([]string, 0),
	}
	err = parse(&updateconverter, s)
	if err != nil {
		return
	}
	sql, data, err = updateconverter.toSql()
	return
}

func Convert2DeleteSql(s interface{}) (sql string, data []interface{}, err error) {
	deleteconverter := deleteconverter{
		conditions: make([]squirrel.Sqlizer, 0),
		order:      make([]string, 0),
	}
	err = parse(&deleteconverter, s)
	if err != nil {
		return
	}
	sql, data, err = deleteconverter.toSql()
	return
}
