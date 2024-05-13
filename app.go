package main

import (
	"consulate/routes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/enquiries", routes.StoreEnquiry)
	r.POST("/exotel/callback", routes.ExotelCallback)
	r.POST("/slack/interaction", routes.SlackInteraction)
	r.POST("/slack/options", routes.SlackOptions)
	r.Run()
}
