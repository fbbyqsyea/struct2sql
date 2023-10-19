package struct2sql

func Convert2SelectSql() {}

func Convert2InsertSql(s interface{}) (sql string, data []interface{}, err error) {
	insertconverter := insertconverter{
		columns: make([]string, 0),
		values:  make([]interface{}, 0),
	}
	err = parse(&insertconverter, s)
	if err != nil {
		return
	}
	return insertconverter.toSql()
}

func Convert2UpdateSql() {}

func Convert2DeleteSql() {}
