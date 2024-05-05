package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/model"
)

func RobotsTxt(c *gin.Context) {
	c.String(200, "User-Agent: *\nDisallow: /\n")
}

func Unsubscribe(c *gin.Context) {
	templates := config.Content.WebTemplate
	c.Header("Content-Type", "text/html; charset=utf-8")
	post, thread, err := model.ConfirmUnsubscribe(c.Param("postId"), c.Param("threadId"), c.Param("editKey"))
	if err != nil {
		fmt.Printf("Unable to prepare unsubscribe: %s\n", err)
		c.String(404, templates.UnsubscribeError.Render())
		return
	}
	templateData := map[string]interface{}{
		"Post":   *post,
		"Thread": *thread,
	}
	c.String(200, templates.UnsubscribeConfirm.Render(templateData))
}

func UnsubscribeConfirm(c *gin.Context) {
	templates := config.Content.WebTemplate
	confirmed := c.PostForm("userConfirmed") == "true"
	c.Header("Content-Type", "text/html; charset=utf-8")
	post, thread, err := model.PerformUnsubscribe(c.Param("postId"), c.Param("threadId"), c.Param("editKey"), confirmed)
	if err != nil {
		fmt.Printf("Unable to unsubscribe: %s\n", err)
		c.String(500, templates.UnsubscribeError.Render())
		return
	}

	templateData := map[string]interface{}{
		"Post":   *post,
		"Thread": *thread,
	}
	c.String(200, templates.UnsubscribeSuccess.Render(templateData))
}
