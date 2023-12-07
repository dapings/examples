package api

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/service"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/convert"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/upload"
	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(ctx *gin.Context) {
	res := app.NewRes(ctx)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		res.ToErrRes(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(ctx.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		res.ToErrRes(errcode.InvalidParams)
		return
	}

	svc := service.New(ctx.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(ctx, "service.UploadFile err: %v", err)
		res.ToErrRes(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	res.ToRes(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
