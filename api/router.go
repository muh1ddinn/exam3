package api

import (
	"exam3/api/handler"
	"exam3/pkg/logger"
	"exam3/service"

	"github.com/gin-gonic/gin"

	_ "exam3/api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceMangaer, log logger.ILogger) *gin.Engine {
	h := handler.NewStrg(services, log)

	r := gin.Default()

	r.POST("/user", h.Createcus)
	r.PUT("/user", h.Update)
	r.GET("/user/:id", h.GetByID)
	r.GET("/user", h.Getalluser)
	r.DELETE("/user/:id", h.Deleteuser)
	r.PUT("/user/password", h.Changepassword)
	r.PATCH("/status", h.Updatestatus)

	r.POST("/users/login", h.UserRegister)
	r.POST("/users/register", h.UserRegister)
	r.POST("/users/register-confirm", h.UsersRegisterConfirm)
	r.POST("/users/login_for_otp", h.Userloginwith_otp)
	r.POST("/users/login_with_otp", h.UsersLoginwithotp)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
