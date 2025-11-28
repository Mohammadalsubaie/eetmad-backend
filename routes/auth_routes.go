package routes

import (
	"github.com/eetmad/backend/controllers"
	"github.com/eetmad/backend/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login",    controllers.Login)
		auth.POST("/refresh",  controllers.Refresh)
	}

	// الروت الصحيح والأخير إن شاء الله
	protected := r.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/api/v1/user/me", controllers.Me)
	}
}
