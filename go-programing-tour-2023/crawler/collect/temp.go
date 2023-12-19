package collect

type Temp struct {
	data map[string]any
}

// Get 返回临时缓存数据
func (t *Temp) Get(key string) any {
	return t.data[key]
}

func (t *Temp) Set(key string, val any) error {
	if t.data == nil {
		t.data = make(map[string]any)
	}

	t.data[key] = val

	return nil
}
