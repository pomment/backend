package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pomment/pomment/auth"
	"github.com/pomment/pomment/utils"
	"net/http"
	"strings"
)

func VerifyToken(c *gin.Context) {
	reqToken := c.GetHeader("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		utils.AjaxError(c, http.StatusForbidden, errors.New("admin token is invalid"))
		c.Abort()
		return
	}
	reqToken = splitToken[1]
	_, err := auth.ValidateToken(reqToken)
	if err != nil {
		println("Admin token validation failed!")
		utils.AjaxError(c, http.StatusForbidden, err)
		c.Abort()
		return
	}
	c.Next()
}

func NoVerifyToken(c *gin.Context) {
	c.Next()
}
