package middleware

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		if s, exist := ctx.GetQuery("token"); exist {
			token = s
		} else {
			token = ctx.GetHeader("token")
		}
		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}

		if ecode != errcode.Success {
			resp := app.NewRes(ctx)
			resp.ToErrRes(ecode)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
