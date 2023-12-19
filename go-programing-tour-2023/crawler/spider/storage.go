package spider

type (
	DataCell struct {
		Task *Task
		Data map[string]interface{}
	}

	Storage interface {
		Save(datas ...*DataCell) error
	}
)

func (d *DataCell) GetTaskName() string {
	return d.Data["Task"].(string)
}

func (d *DataCell) GetTableName() string {
	return d.Data["Task"].(string)
}
