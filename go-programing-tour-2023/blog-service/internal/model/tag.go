package model

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"status"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}
