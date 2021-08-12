package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/wolf88804/blog-service/docs"
	"github.com/wolf88804/blog-service/global"
	"github.com/wolf88804/blog-service/internal/middlware"
	api "github.com/wolf88804/blog-service/internal/routers/api"
	v1 "github.com/wolf88804/blog-service/internal/routers/api/v1"
	"github.com/wolf88804/blog-service/pkg/limiter"
	"net/http"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "info" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middlware.AccessLog())
		r.Use(middlware.Recovery())
	}
	r.Use(middlware.AppInfo())
	//限流
	r.Use(middlware.RateLimiter(methodLimiters))
	//限时
	r.Use(middlware.ContextTimeout(60 * time.Second))
	//追踪
	r.Use(middlware.Tracing())
	//中间件校验中文
	r.Use(middlware.Translations())
	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//上传文件
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	//认证
	r.GET("/auth", api.GetAuth)

	tag := v1.NewTag()
	article := v1.NewArticle()
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middlware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags/:id", tag.Get)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
