package routes

import (
	"consulate/helpers"
	"consulate/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StoreEnquiry(c *gin.Context) {
	var enquiry models.Enquiry
	if err := c.ShouldBind(&enquiry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enquiry.Status = "fresh"
	if err := helpers.HandleNewEnquiry(enquiry, c.ClientIP()); err != nil {
		panic(err)
	}
}
