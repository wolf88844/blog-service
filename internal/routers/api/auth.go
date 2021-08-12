package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wolf88804/blog-service/global"
	"github.com/wolf88804/blog-service/internal/service"
	"github.com/wolf88804/blog-service/pkg/app"
	"github.com/wolf88804/blog-service/pkg/errcode"
)

// @Summary 获取token
// @Produce json
// @Param app_key query string true "Key"
// @Param app_secret query string true "Secret"
// @Success 200 {object} gin.Context "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	response := app.NewResponse(c)
	vaild, errs := app.BindAndValid(c, &param)
	if vaild {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf(c,"svc.CheckAuth err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}
	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c,"svc.GenerateToken err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}
	response.ToResponse(gin.H{"token": token})
}
