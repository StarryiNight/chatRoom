package routes

import (
	"chatRoom/controllers"
	"chatRoom/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	//注册业务路由
	r.POST("/signup", controllers.SignUpHandler)
	//登陆业务路由
	r.POST("/login", controllers.LoginHandler)

	r.GET("/", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/logout", middlewares.JWTAuthMiddleware(), controllers.Logout)
	r.GET("/chat", controllers.Server)
	return r
}
