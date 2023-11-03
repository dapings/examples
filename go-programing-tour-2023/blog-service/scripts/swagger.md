# OpenAPI - Swagger

Swagger 是基于标准的 OpenAPI 规范进行设计的，只要按照这套规范编写注解或通过扫描代码生成注解，就能生成统一标准的接口文档和一系列 Swagger 工具。

## 安装 Swagger

```shell
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/swaggo/gin-swagger@latest
#github.com/swaggo/files@v2.0.0
#github.com/alecthomas/template
```
```shell
swag -v
```

## 写入注解

| 注解       | 描述                                      |
|----------|-----------------------------------------|
| @Summary | 摘要                                      |
| @Produce | 响应类型，如 JSON, XML, HTML等                 |
| @Param   | 参数格式：参数名、参数类型、数据类型、是否必填、注释              |
| @Success | 响应成功：状态码、参数类型、数据类型、注释                   |
| @Failure | 响应失败：状态码、参数类型、数据类型、注释                   |
| @Router  | 路由：路由地址、HTTP方法，e.g. /api/v1/debug [get] |

```text
// @Summary 创建文章
// @Produce json
// @Param tag_id body string true "标签ID"
// @Param title body string true "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string true "封面图片地址"
// @Param content body string true "文章内容"
// @Param created_by body int true "创建者"
// @Param state body int false "状态"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [post]
```

## main 入口添加注解

```text
// @title 博客系统
// @version 1.0
// @description Go 语言 blog service
// @termsOfService https://github.com/dapings/examples
```

## 生成

```shell
# default main.go
swag init
swag init -g blog.go
```

在 docs 目录下生成了 docs.go, swagger.json, swagger.yaml 三个文件。

## 访问接口文档

在路由中进行默认初始化和注册对就的路由：
```go
package main
import (
	"github.com/gin-gonic/gin"
	// 初始化 docs 包时，会执行 docs/docs.go 的 init 方法，关联文档
    // 接着在 ReadDoc 方法中做一些 template 的模板映射等工作
    _ "github.com/dapings/examples/go-programing-tour-2023/blog-service/docs"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
	
    // 将调用的URL设置为 doc.json(生成的 swagger 注解)
	// 方式一
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    // 方式二
	// swagger.josn 会默认指向当前应用所启动的域名下的 swagger/doc.json 路径，如有额外需求，可手动指定：
	url := ginSwagger.URL("http://0.0.0.0:8000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	
	return r
}
```

## 查看接口文档

访问 Swagger 的地址：http://127.0.0.1:8000/swagger/index.html
主要分为三部分：项目主体信息、接口路由、模型信息。

## 当有 model.Tag 以外的字段时，如分页等，如何展示

官方建议，定义一个针对 Swagger 的对象，专门用于Swagger接口文档的展示：
```
type ArticleSwagger struct {
	...
}
```
```
// @Success 200 {object} model.ArticleSwagger "成功"
```
调整接口方法中对应的注解信息。