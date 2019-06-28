package v1

import (
	"github.com/gin-gonic/gin"
	"github_statistics/internal/log"
	"github_statistics/pkg/github"
	"io/ioutil"
	"net/http"
)

func Router(engine *gin.Engine) {
	r := engine.Group("/api/v1")

	r.POST("/webhook/callback", WebhookCallback)
}

func WebhookCallback(c *gin.Context) {
	log.Info("Header: ", c.Request.Header)
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = github.HandleEvent(c.GetHeader("X-GitHub-Event"), data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

}