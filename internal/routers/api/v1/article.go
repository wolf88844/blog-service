package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wolf88804/blog-service/global"
	"github.com/wolf88804/blog-service/internal/service"
	"github.com/wolf88804/blog-service/pkg/app"
	"github.com/wolf88804/blog-service/pkg/convert"
	"github.com/wolf88804/blog-service/pkg/errcode"
)

type Article struct {
}

func NewArticle() Article {
	return Article{}
}

// @Summary 获取单个文章
// @Produce json
// @Param token header string true "token"
// @Param id path int true "文章ID"
// @Param state query int false "状态" Enums(0,1) default(1)
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [get]
func (a *Article) Get(c *gin.Context) {
	param := service.ArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if valid {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	article, err := svc.GetArticle(&param)
	if err != nil {
		global.Logger.Errorf(c,"svc.GetArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetArticleFail)
		return
	}
	response.ToResponse(article)
	return
}

// @Summary 获取多个个文章
// @Produce json
// @Param tag_id query int false "标签ID"
// @Param title query string false "文章标题" maxlength(100)
// @Param state query int false "状态" Enums(0,1) default(1)
// @Param page query int false "页码"
// @Prams page_size query int false "每页数量"
// @Success 200 {object} model.ArticleSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [get]
func (a *Article) List(c *gin.Context) {
	param := service.ArticleListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if valid {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	articles, totalRows, err := svc.GetArticleList(&param, &pager)
	if err != nil {
		global.Logger.Errorf(c,"svc.GetArticleList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetArticleListFail)
		return
	}
	response.ToResponseList(articles, totalRows)
	return
}

// @Summary 新增文章
// @Produce json
// @Param title query string true "文章标题" maxlength(100) minlength(3)
// @Param state query int false "状态" Enums(0,1) default(1)
// @Param desc query string false "文章简述" minlength(3) maxlength(100)
// @Param cover_image_url query string false "封图片地址" minlength(3) maxlength(100)
// @Param content query string true "文章内容" minlength(3)
// @Param created_by query string false "创建者" minlength(3) maxlength(100)
// @Param tag_id query int false "标签ID"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [post]
func (a *Article) Create(c *gin.Context) {
	param := service.CreateArticleRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if valid {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CreateArticle(&param)
	if err != nil {
		global.Logger.Errorf(c,"svc.CreateArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateArticleFail)
		return
	}
	response.ToResponse(gin.H{})
	return
}

// @Summary 更新文章
// @Produce json
// @Param id path int true "文章ID"
// @Param title query string false "文章标题" maxlength(100) minlength(3)
// @Param state query int false "状态" Enums(0,1) default(1)
// @Param modified_by query string true "修改者" minlength(3) maxlength(100)
// @Param desc query string false "文章简述" minlength(3) maxlength(100)
// @Param cover_image_url query string false "封图片地址" minlength(3) maxlength(100)
// @Param content query string false "文章内容" minlength(3)
// @Param tag_id query int false "标签ID"
// @Success 200 {array} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [put]
func (a *Article) Update(c *gin.Context) {
	param := service.UpdateArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if valid {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.UpdateArticle(&param)
	if err != nil {
		global.Logger.Errorf(c,"svc.UpdateArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateArticleFail)
		return
	}
	response.ToResponse(gin.H{})
	return
}

// @Summary 删除文章
// @Produce json
// @Param id path int true "文章ID"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [delete]
func (a *Article) Delete(c *gin.Context) {
	param := service.DeleteArticleRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if valid {
		global.Logger.Errorf(c,"app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.DeleteArticle(&param)
	if err != nil {
		global.Logger.Errorf(c,"svc.DeleteArticle err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteArticleFail)
		return
	}
	response.ToResponse(gin.H{})
	return
}
