package global

import (
	"github.com/wolf88804/blog-service/pkg/logger"
	"github.com/wolf88804/blog-service/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
	JWTSetting      *setting.JWTSetting
	EmailSetting    *setting.EmailSetting
)
