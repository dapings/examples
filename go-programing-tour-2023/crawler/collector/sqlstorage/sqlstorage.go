package sqlstorage

import (
	"encoding/json"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/collector"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/engine"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"go.uber.org/zap"
)

type SQLStore struct {
	dataCells   []*collector.DataCell // 分批输出结果缓存
	columnNames []sqldb.Field         // 标题字段
	db          sqldb.DBer
	Table       map[string]struct{}
	options
}

func New(opts ...Option) (*SQLStore, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &SQLStore{}
	s.options = options
	s.Table = make(map[string]struct{})
	var err error
	s.db, err = sqldb.New(sqldb.WithConnUrl(s.sqlUrl), sqldb.WithLogger(s.logger))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLStore) Save(dataCells ...*collector.DataCell) error {
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
				return err
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
func getFields(cell *collector.DataCell) []sqldb.Field {
	taskName := cell.Data["Task"].(string)
	ruleName := cell.Data["Rule"].(string)
	fields := engine.GetFields(taskName, ruleName)

	var columnNames []sqldb.Field
	for _, field := range fields {
		columnNames = append(columnNames, sqldb.Field{Title: field, Type: "MEDIUMTEXT"})
	}
	columnNames = append(columnNames,
		sqldb.Field{Title: "Url", Type: "VARCHAR(255)"},
		sqldb.Field{Title: "Time", Type: "VARCHAR(255)"},
	)
	return columnNames
}

func (s *SQLStore) Flush() error {
	if len(s.dataCells) == 0 {
		return nil
	}

	args := make([]any, 0)
	for _, dataCell := range s.dataCells {
		taskName := dataCell.Data["Task"].(string)
		ruleName := dataCell.Data["Rule"].(string)
		fields := engine.GetFields(taskName, ruleName)
		data := dataCell.Data["Data"].(map[string]any)
		var vals []string
		for _, field := range fields {
			v := data[field]
			switch v.(type) {
			case nil:
				vals = append(vals, "")
			case string:
				vals = append(vals, v.(string))
			default:
				buf, err := json.Marshal(v)
				if err != nil {
					vals = append(vals, "")
				} else {
					vals = append(vals, string(buf))
				}

			}
		}
		vals = append(vals, dataCell.Data["Url"].(string), dataCell.Data["Time"].(string))
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
