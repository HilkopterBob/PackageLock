package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Data": "hello from controller!"})
}
