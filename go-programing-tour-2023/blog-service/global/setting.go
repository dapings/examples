package global

import (
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/logger"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	EmailSetting    *setting.EmailSettingS
	JWTSetting      *setting.JWTSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
