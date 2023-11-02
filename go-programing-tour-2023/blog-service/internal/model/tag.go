package model

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"status"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}
