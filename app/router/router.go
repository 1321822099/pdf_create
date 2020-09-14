package router

import (
	"fmt"

	"github.com/1321822099/pdf_create/app/controllers"
	"github.com/1321822099/pdf_create/app/utils/config"
	"github.com/gin-gonic/gin"
)

// Run router
func Run() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello World")
    })
    router.Run(":8000")
}
