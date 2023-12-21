package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	tableName := "test_create_table"
	var notValidTable = TableData{
		TableName: tableName,
		ColumnNames: []Field{
			{Title: "书名", Type: "not_valid"},
			{Title: "URL", Type: "VARCHAR(255)"},
		},
		AutoKey: true,
	}

	sqldb, err := New(WithConnURL(ConnStrWithMySQL))
	assert.Nil(t, err)
	assert.NotNil(t, sqldb)

	defer func() {
		err := sqldb.DropTable(notValidTable)
		assert.Nil(t, err)
	}()

	err = sqldb.CreateTable(notValidTable)
	assert.Nil(t, err)

	var validTable = TableData{
		TableName: tableName,
		ColumnNames: []Field{
			{Title: "书名", Type: "MEDIUMTEXT"},
			{Title: "URL", Type: "VARCHAR(255)"},
		},
		AutoKey: true,
	}
	err = sqldb.CreateTable(validTable)
	assert.NotNil(t, err)
}

func TestCreateTableByTableDriver(t *testing.T) {
	tableName := "test_create_table"

	type args struct {
		t TableData
	}

	testDatas := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid_table",
			args: args{TableData{
				TableName: tableName,
				ColumnNames: []Field{
					{Title: "书名", Type: "not_valid"},
					{Title: "URL", Type: "VARCHAR(255)"},
				},
			}},
			wantErr: true,
		},
		{
			name: "valid_table",
			args: args{TableData{
				TableName: tableName,
				ColumnNames: []Field{
					{Title: "书名", Type: "MEDIUMTEXT"},
					{Title: "URL", Type: "VARCHAR(255)"},
				},
			}},
			wantErr: false,
		},
		{
			name: "valid_table_with_pk",
			args: args{TableData{
				TableName: tableName,
				ColumnNames: []Field{
					{Title: "书名", Type: "MEDIUMTEXT"},
					{Title: "URL", Type: "VARCHAR(255)"},
				},
				AutoKey: true,
			}},
			wantErr: false,
		},
	}

	sqldb, err := New(WithConnURL(ConnStrWithMySQL))
	assert.Nil(t, err)
	assert.NotNil(t, sqldb)

	for _, tt := range testDatas {
		err = sqldb.CreateTable(tt.args.t)
		if tt.wantErr {
			assert.NotNil(t, err, tt.name)
		} else {
			assert.Nil(t, err, tt.name)
		}

		err = sqldb.DropTable(tt.args.t)
		assert.Nil(t, err, tt.name)
	}
}

func TestInsert(t *testing.T) {
	tableName := "test_create_table"
	columnNames := []Field{{Title: "书名", Type: "MEDIUMTEXT"}, {Title: "price", Type: "TINYINT"}}

	type args struct {
		t TableData
	}

	testDatas := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "insert_one_data",
			args: args{TableData{
				TableName:   tableName,
				ColumnNames: columnNames,
				Args:        []any{"book1", 2},
				DataCount:   1,
			}},
			wantErr: false,
		},
		{
			name: "insert_multi_data",
			args: args{TableData{
				TableName:   tableName,
				ColumnNames: columnNames,
				Args:        []any{"book2", 79, "book3", 66.88},
				DataCount:   2,
			}},
			wantErr: false,
		},
		{
			name: "insert_multi_data_but_wrong_count",
			args: args{TableData{
				TableName:   tableName,
				ColumnNames: columnNames,
				Args:        []any{"book2", 79, "book3", 66.88},
				DataCount:   1,
			}},
			wantErr: true,
		},
		{
			name: "insert_wrong_data_type",
			args: args{TableData{
				TableName:   tableName,
				ColumnNames: columnNames,
				Args:        []any{"book2", "none"},
				DataCount:   1,
			}},
			wantErr: false,
		},
	}

	sqldb, err := New(WithConnURL(ConnStrWithMySQL))
	assert.Nil(t, err)
	assert.NotNil(t, sqldb)

	err = sqldb.CreateTable(testDatas[0].args.t)
	assert.Nil(t, err, testDatas[0].name)

	defer func() {
		err = sqldb.DropTable(testDatas[0].args.t)
		assert.Nil(t, err, testDatas[0].name)
	}()

	for _, tt := range testDatas {
		t.Run(tt.name, func(t *testing.T) {
			err = sqldb.Insert(tt.args.t)
			if tt.wantErr {
				assert.NotNil(t, err, tt.name)
			} else {
				assert.Nil(t, err, tt.name)
			}
		})
	}
}
