package v1

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Article struct{}

func NewArticle() Article {
	return Article{}
}

func (a Article) Get(ctx *gin.Context) {
	app.NewRes(ctx).ToErrRes(errcode.ServerError)
	return
}

func (a Article) List(ctx *gin.Context) {

}

func (a Article) Create(ctx *gin.Context) {

}

func (a Article) Update(ctx *gin.Context) {

}

func (a Article) Delete(ctx *gin.Context) {

}
