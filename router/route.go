package router

import (
	"srvguide/internal/app"

	"github.com/gin-gonic/gin"
)

//Init 初始化路由信息
func Init(r *gin.Engine) {
	apiGroup := r.Group("/rest/v1")
	apiGroup.GET("/racing", app.GetFunction)
}
