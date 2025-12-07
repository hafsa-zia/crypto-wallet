package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONError is a small helper for error responses.
func JSONError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// JSONOK is a helper for success responses.
func JSONOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
