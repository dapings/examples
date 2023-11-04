package dao

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/model"
)

func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
	auth := model.Auth{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
	return auth.Get(d.engine)
}
