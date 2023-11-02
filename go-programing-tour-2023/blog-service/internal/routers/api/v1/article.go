package v1

import (
	"github.com/gin-gonic/gin"
)

type Article struct{}

func NewArticle() Article {
	return Article{}
}

func (a Article) Get(ctx *gin.Context) {

}

func (a Article) List(ctx *gin.Context) {

}

func (a Article) Create(ctx *gin.Context) {

}

func (a Article) Update(ctx *gin.Context) {

}

func (a Article) Delete(ctx *gin.Context) {

}
