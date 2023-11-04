package routers

import (
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/middleware"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers/api"
	v1 "github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Translations())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"message": "pong"})
	})

	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	// 注：应当将文件服务和应用服务分开，因为从安全角度，文件资源不应当于应用资源在一起，或直接使用 oss 也是可以的。
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	// 获取 JWT Token
	r.GET("/auth", api.GetAuth)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWT())
	{
		tag := v1.NewTag()
		apiV1.POST("/tags", tag.Create)
		apiV1.DELETE("/tags/:id", tag.Delete)
		apiV1.PUT("/tags/:id", tag.Update)
		apiV1.PATCH("/tags/:id/state", tag.Update)
		apiV1.GET("/tags", tag.List)

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
