package sql2struct

import (
	"fmt"
	"os"
	"text/template"

	"github.com/dapings/examples/go-programing-tour-2023/tour/internal/word"
)

const structTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`

var tplName = "sql2struct"

type ConvStructTpl struct {
	convTpl string
}

type ConvStructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type ConvStructTplDB struct {
	TableName string
	Columns   []*ConvStructColumn
}

func NewConvStructTpl() *ConvStructTpl {
	return &ConvStructTpl{
		convTpl: structTpl,
	}
}

func (t *ConvStructTpl) AssemblyColumns(tableColumns []*TableColumn) []*ConvStructColumn {
	tplColumns := make([]*ConvStructColumn, 0, len(tableColumns))
	for _, column := range tableColumns {
		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", column.ColumnName)
		tplColumns = append(tplColumns, &ConvStructColumn{
			Name:    column.ColumnName,
			Type:    column.ColumnType,
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *ConvStructTpl) Generate(tableName string, tplColumns []*ConvStructColumn) error {
	tpl := template.Must(template.New(tplName).Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.convTpl))
	tplDB := ConvStructTplDB{
		TableName: tableName,
		Columns:   tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}
	return nil
}
