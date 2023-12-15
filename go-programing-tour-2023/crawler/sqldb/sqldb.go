package sqldb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type SQLDB struct {
	options
	db *sql.DB
}

type Field struct {
	Title string
	Type  string
}

type TableData struct {
	TableName   string
	ColumnNames []Field // 标题字段
	Args        []any   // 数据
	DataCount   int     // 插入数据的量
	AutoKey     bool
}
