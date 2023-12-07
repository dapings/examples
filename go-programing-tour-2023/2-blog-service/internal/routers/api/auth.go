package api

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/service"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func GetAuth(ctx *gin.Context) {
	param := service.AuthRequest{}
	resp := app.NewRes(ctx)
	valid, errs := app.BindAndValid(ctx, &param)
	if !valid {
		global.Logger.Errorf(ctx, "app.BindAndValid errs: %v", errs)
		resp.ToErrRes(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(ctx.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf(ctx, "service.CheckAuth err: %v", err)
		resp.ToErrRes(errcode.UnauthorizedAuthNotExist)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(ctx, "app.GenerateToken err: %v", err)
		resp.ToErrRes(errcode.UnauthorizedTokenGenerate)
		return
	}

	resp.ToRes(gin.H{
		"token": token,
	})
}
