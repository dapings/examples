package sqlstorage

import (
	"encoding/json"
	"errors"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"go.uber.org/zap"
)

type SQLStorage struct {
	dataCells []*spider.DataCell // 分批输出结果缓存
	// columnNames []sqldb.Field       // 标题字段
	db    sqldb.DBer
	Table map[string]struct{}
	options
}

func New(opts ...Option) (*SQLStorage, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	s := &SQLStorage{}
	s.options = options
	s.Table = make(map[string]struct{})

	var err error
	s.db, err = sqldb.New(sqldb.WithConnURL(s.sqlURL), sqldb.WithLogger(s.logger))

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLStorage) Save(dataCells ...*spider.DataCell) error {
	for _, cell := range dataCells {
		name := cell.GetTableName()
		if _, ok := s.Table[name]; !ok {
			columnNames := getFields(cell)
			err := s.db.CreateTable(sqldb.TableData{
				TableName:   name,
				ColumnNames: columnNames,
				AutoKey:     true,
			})

			if err != nil {
				s.logger.Error("create table failed", zap.Error(err))

				continue
			}

			s.Table[name] = struct{}{}
		}

		if len(s.dataCells) >= s.BatchCount {
			err := s.Flush()

			if err != nil {
				s.logger.Error("flush failed", zap.Error(err))

				return err
			}
		}

		s.dataCells = append(s.dataCells, cell)
	}

	return nil
}

// 解析出标题字段
func getFields(cell *spider.DataCell) []sqldb.Field {
	taskName := cell.Data["Task"].(string)
	ruleName := cell.Data["Rule"].(string)
	fields := engine.GetFields(taskName, ruleName)

	var columnNames []sqldb.Field
	for _, field := range fields {
		columnNames = append(columnNames, sqldb.Field{Title: field, Type: "MEDIUMTEXT"})
	}

	columnNames = append(columnNames,
		sqldb.Field{Title: "URL", Type: "VARCHAR(255)"},
		sqldb.Field{Title: "Time", Type: "VARCHAR(255)"},
	)

	return columnNames
}

func (s *SQLStorage) Flush() error {
	if len(s.dataCells) == 0 {
		return nil
	}

	defer func() {
		s.dataCells = nil
	}()

	args := make([]any, 0)

	var taskName, ruleName string
	var ok bool
	for _, dataCell := range s.dataCells {
		if ruleName, ok = dataCell.Data["Rule"].(string); !ok {
			return errors.New("no rule field")
		}

		if taskName, ok = dataCell.Data["Task"].(string); !ok {
			return errors.New("no task field")
		}

		fields := engine.GetFields(taskName, ruleName)

		data := dataCell.Data["Data"].(map[string]any)

		var vals []string

		for _, field := range fields {
			v := data[field]
			switch v := v.(type) {
			case nil:
				vals = append(vals, "")
			case string:
				vals = append(vals, v)
			default:
				buf, err := json.Marshal(v)
				if err != nil {
					vals = append(vals, "")
				} else {
					vals = append(vals, string(buf))
				}
			}
		}

		if v, ok := dataCell.Data["URL"].(string); ok {
			vals = append(vals, v)
		}

		if v, ok := dataCell.Data["Time"].(string); ok {
			vals = append(vals, v)
		}

		for _, val := range vals {
			args = append(args, val)
		}
	}

	return s.db.Insert(sqldb.TableData{
		TableName:   s.dataCells[0].GetTableName(),
		ColumnNames: getFields(s.dataCells[0]),
		Args:        args,
		DataCount:   len(s.dataCells),
	})
}
