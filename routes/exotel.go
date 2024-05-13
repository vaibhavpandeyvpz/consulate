package routes

import (
	"consulate/helpers"
	"consulate/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExotelCallback(c *gin.Context) {
	var call models.ExotelCall
	if err := c.ShouldBind(&call); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if call.CustomField == "" || call.EventType != "terminal" || call.Status != "completed" {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := helpers.SendCallRecording(call); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
