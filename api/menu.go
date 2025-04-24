package api

import "github.com/gin-gonic/gin"

func InitSecuredMenuRouter(r *gin.RouterGroup) {
	securedRouter := r.Group("/auth")
	securedRouter.POST("/add", addMenu)
}

func addMenu(c *gin.Context) {

}
