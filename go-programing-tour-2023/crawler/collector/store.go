package collector

type (
	OutputData struct {
		Data map[string]interface{}
	}

	Store interface {
		Save(datas ...OutputData) error
	}
)
