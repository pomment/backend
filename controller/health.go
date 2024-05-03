package controller

import (
	"github.com/gin-gonic/gin"
	"pomment-go/utils"
)

func Health(c *gin.Context) {
	utils.AjaxSuccess(c, nil)
}
