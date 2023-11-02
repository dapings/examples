package routers

import (
	"net/http"

	v1 "github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"message": "pong"})
	})

	apiV1 := r.Group("/api/v1")
	{
		tag := v1.NewTag()
		apiV1.POST("/tags", tag.Create)
		apiV1.DELETE("/tags/:id", tag.Delete)
		apiV1.PUT("/tags/:id", tag.Update)
		apiV1.PATCH("/tags/:id/state", tag.Update)
		apiV1.GET("/tags", tag.Get)

		article := v1.NewArticle()
		apiV1.POST("/articles", article.Create)
		apiV1.DELETE("/articles/:id", article.Delete)
		apiV1.PUT("/articles/:id", article.Update)
		apiV1.PATCH("/articles/:id/state", article.Update)
		apiV1.GET("/articles/:id", article.Get)
		apiV1.GET("/articles", article.Get)
	}

	return r
}