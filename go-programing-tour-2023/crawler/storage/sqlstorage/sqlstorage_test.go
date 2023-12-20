package sqlstorage

import (
	"testing"

	"github.com/dapings/examples/go-programing-tour-2023/crawler/spider"
	"github.com/dapings/examples/go-programing-tour-2023/crawler/sqldb"
	"github.com/stretchr/testify/assert"
)

type mysqldb struct{}

func (m mysqldb) CreateTable(_ sqldb.TableData) error {
	return nil
}

func (m mysqldb) Insert(_ sqldb.TableData) error {
	return nil
}

func TestFlush(t *testing.T) {
	type fields struct {
		dataCells []*spider.DataCell
		options   options
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "empty", wantErr: false},
		{name: "no rule filed", fields: fields{dataCells: []*spider.DataCell{{Data: map[string]interface{}{"url": "https://x.cn"}}}},
			wantErr: true},
		{name: "no task filed", fields: fields{dataCells: []*spider.DataCell{{Data: map[string]interface{}{"url": "https://x.cn"}}}},
			wantErr: true},
		{name: "right data", fields: fields{dataCells: []*spider.DataCell{
			{Data: map[string]interface{}{"Rule": "书籍简介", "Task": "douban_book_list", "Data": map[string]interface{}{"url": "https://x.cn"}}},
		}}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLStorage{
				dataCells: tt.fields.dataCells,
				db:        mysqldb{},
				options:   tt.fields.options,
			}
			if err := s.Flush(); (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, but want error %v", err, tt.wantErr)
			}

			assert.Nil(t, s.dataCells)
		})
	}
}
