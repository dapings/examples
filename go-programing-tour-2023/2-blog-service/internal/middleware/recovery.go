package middleware

import (
	"fmt"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/app"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/email"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(ctx, "panic recover err: %v", err)

				err := defaultMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("exception happen: %d", time.Now().Unix()),
					fmt.Sprintf("err: %v", err),
				)
				if err != nil {
					global.Logger.Panicf(ctx, "mail.SendMail err: %v", err)
				}

				app.NewRes(ctx).ToErrRes(errcode.ServerError)
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
