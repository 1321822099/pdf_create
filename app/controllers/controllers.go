package controllers

import (
	"net/http"

	"github.com/1321822099/pdf_create/service/cmd"
	"github.com/gin-gonic/gin"
)

func Command(c *gin.Context) {
	req := &cmd.CommandReq{}
	if err := ShouldBindJSON(c, req); err != nil {
		return
	}
	resp, err := cmd.RunCommand(req)
	RenderJSON(c, err, resp)
}

func ShouldBindJSON(c *gin.Context, req interface{}) error {
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	return err
}

func RenderJSON(c *gin.Context, result error, obj interface{}) {
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error()})
		return
	}
	if obj == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, obj)
}
