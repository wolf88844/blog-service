package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wolf88804/blog-service/global"
	"github.com/wolf88804/blog-service/internal/service"
	"github.com/wolf88804/blog-service/pkg/app"
	"github.com/wolf88804/blog-service/pkg/convert"
	"github.com/wolf88804/blog-service/pkg/errcode"
	"github.com/wolf88804/blog-service/pkg/upload"
)

type Upload struct {
}

func NewUpload() Upload {
	return Upload{}
}

// @Summary 上传文件
// @Produce json
// @Accept multipart/form-data
// @Param file formData file true "文件"
// @Param type formData string true "文件类型"
// @Success 200 {object}  gin.Context "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /upload/file [post]
func (u *Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c,"svc.uploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}
	response.ToResponse(gin.H{"file_access_url": fileInfo.AccessUrl})

}
