package router

import "github.com/gin-gonic/gin"

func UseRegisteredApis(engine *gin.Engine) {
	engine.GET("/echo", func(context *gin.Context) {
		str := context.Query("str")
		context.JSON(200, "echo :"+str)
	})
}
