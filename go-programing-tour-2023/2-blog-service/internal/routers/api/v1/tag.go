package v1

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/service"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/convert"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Tag struct{}

func NewTag() Tag {
	return Tag{}
}

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(ctx *gin.Context) {
	param := service.TagListRequest{}
	resp := app.NewRes(ctx)
	valid, errs := app.BindAndValid(ctx, &param)
	if !valid {
		global.Logger.Errorf(ctx, "app.BindAndValid errs: %v", errs)
		resp.ToErrRes(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(ctx.Request.Context())
	pager := app.Pager{Page: app.GetPage(ctx), PageSize: app.GetPageSize(ctx)}
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.Errorf(ctx, "service.CountTag err: %v", err)
		resp.ToErrRes(errcode.ErrorCountTagFail)
		return
	}
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errorf(ctx, "service.GetTagList err: %v", err)
		resp.ToErrRes(errcode.ErrorGetTagListFail)
		return
	}

	resp.ToResList(tags, totalRows)
	return
}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string false "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(ctx *gin.Context) {
	param := service.CreateTagRequest{}
	resp := app.NewRes(ctx)
	valid, errs := app.BindAndValid(ctx, &param)
	if !valid {
		global.Logger.Errorf(ctx, "app.BindAndValid errs: %v", errs)
		resp.ToErrRes(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(ctx.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.Errorf(ctx, "service.CreateTag err: %v", err)
		resp.ToErrRes(errcode.ErrorCreateTagFail)
		return
	}

	resp.ToRes(gin.H{})
	return
}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(ctx *gin.Context) {
	param := service.UpdateTagRequest{ID: convert.StrTo(ctx.Param("id")).MustUInt32()}
	resp := app.NewRes(ctx)
	valid, errs := app.BindAndValid(ctx, &param)
	if !valid {
		global.Logger.Errorf(ctx, "app.BindAndValid errs: %v", errs)
		resp.ToErrRes(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(ctx.Request.Context())
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.Errorf(ctx, "service.UpdateTag err: %v", err)
		resp.ToErrRes(errcode.ErrorUpdateTagFail)
		return
	}

	resp.ToRes(gin.H{})
	return
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(ctx *gin.Context) {
	param := service.DeleteTagRequest{ID: convert.StrTo(ctx.Param("id")).MustUInt32()}
	resp := app.NewRes(ctx)
	valid, errs := app.BindAndValid(ctx, &param)
	if !valid {
		global.Logger.Errorf(ctx, "app.BindAndValid errs: %v", errs)
		resp.ToErrRes(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(ctx.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errorf(ctx, "service.DeleteTag err: %v", err)
		resp.ToErrRes(errcode.ErrorDeleteTagFail)
		return
	}

	resp.ToRes(gin.H{})
	return
}
