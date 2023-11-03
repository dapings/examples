package app

import (
	"net/http"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

// 响应处理
type Resp struct {
	Ctx *gin.Context
}

type Pager struct {
	// 	页码
	Page int `json:"page"`
	// 	每页数量
	PageSize int `json:"page_size"`
	// 总行数
	TotalRows int `json:"total_rows"`
}

func NewRes(ctx *gin.Context) *Resp {
	return &Resp{
		Ctx: ctx,
	}
}

func (r *Resp) ToRes(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Resp) ToResList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (r *Resp) ToErrRes(err *errcode.Error) {
	res := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		res["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), res)
}
