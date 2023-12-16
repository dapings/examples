package sqldb

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var (
	driverNameWithMySQL = "mysql"
	ConStrWithMySQL     = "root:123456@tcp(127.0.0.1:3326)/crawler?charset=utf8"
)

type DBer interface {
	CreateTable(TableData) error
	Insert(TableData) error
}

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

func newSQLDB() *SQLDB {
	return &SQLDB{}
}

func New(opts ...Option) (*SQLDB, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	d := &SQLDB{}
	d.options = options
	if err := d.OpenDB(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *SQLDB) OpenDB() error {
	db, err := sql.Open(driverNameWithMySQL, d.sqlUrl)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(2048)
	db.SetMaxIdleConns(2048)
	if err := db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *SQLDB) CreateTable(t TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("when create table, column can not be empty")
	}
	sqlStatement := `CREATE TABLE IF NOT EXISTS ` + t.TableName + ` (`
	if t.AutoKey {
		sqlStatement += `id INT(12) NOT NULL PRIMARY KEY AUTO_INCREMENT,`
	}
	for _, t := range t.ColumnNames {
		sqlStatement += t.Title + ` ` + t.Type + `,`
	}
	sqlStatement = sqlStatement[:len(sqlStatement)-1] + `) ENGINE=InnoDB default CHARSET=utf8mb4;`

	_, err := d.db.Exec(sqlStatement)
	return err
}

func (d *SQLDB) Insert(t TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("when insert data, column can not be empty")
	}
	sqlStatement := `INSERT INTO ` + t.TableName + `(`
	for _, v := range t.ColumnNames {
		sqlStatement += v.Title + `,`
	}
	sqlStatement = sqlStatement[:len(sqlStatement)-1] + `) VALUES `
	blank := `,(` + strings.Repeat(",?", len(t.ColumnNames))[1:] + `)`
	sqlStatement += strings.Repeat(blank, t.DataCount)[1:] + `;`
	d.logger.Debug("insert data", zap.String("sql", sqlStatement))
	_, err := d.db.Exec(sqlStatement, t.Args...)
	return err
}
