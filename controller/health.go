package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pomment/pomment/utils"
)

func Health(c *gin.Context) {
	utils.AjaxSuccess(c, nil)
}
