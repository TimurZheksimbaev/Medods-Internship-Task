package router

import (
	"github.com/gin-gonic/gin"
	"medods-internship/controllers"

)

func Init(tokenController *controllers.TokenController) *gin.Engine {
	router := gin.Default()

	authRouter := router.Group("/auth")
	authRouter.GET("/generate-tokens", tokenController.GenerateTokens)
	authRouter.POST("/refresh-tokens", tokenController.RefreshTokens)
	
	return router
}