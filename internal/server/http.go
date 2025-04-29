package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lrx0014/ScalableFlake/pkg/snowflake"
	log "github.com/sirupsen/logrus"
)

func NewHTTPServer() *gin.Engine {
	r := gin.Default()

	r.GET("/generate_uid", func(c *gin.Context) {
		tenantID := c.Query("tenant_id")
		uid, err := snowflake.GenerateUID(tenantID)
		if err != nil {
			log.Errorf("failed to generate uid: %v", err)
			c.JSON(500, gin.H{"error": "failed to generate UID"})
			return
		}

		c.JSON(200, gin.H{"uid": uid})
	})

	return r
}
