package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xieyuschen/grillen/template/router"
)

func main() {
	r := gin.Default()
	router.UseRegisteredApis(r)
	_ = r.Run(":8080")
}
