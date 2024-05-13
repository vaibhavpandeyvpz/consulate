package routes

import (
	"consulate/helpers"
	"consulate/models"
	"consulate/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SlackInteraction(c *gin.Context) {
	if !services.VerifyRequestFromSlack(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Signature from Slack is not valid."})
		return
	}

	var input models.SlackPayload
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload models.SlackInteractionPayload
	if err := json.Unmarshal([]byte(input.Payload), &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Type == "block_actions" {
		for _, action := range payload.Actions {
			if action.Value == "" {
				continue
			}

			var enquiry models.Enquiry
			id, _ := strconv.Atoi(action.Value)
			services.GormDb().Find(&enquiry, id)

			if action.ID == "call_now" {
				if err := helpers.PlacePhoneCall(payload.User.ID, enquiry); err != nil {
					panic(err)
				}
			} else if action.ID == "follow_up_enquiry" {
				if err := helpers.ShowFollowUpView(payload.TriggerID, enquiry); err != nil {
					panic(err)
				}
			} else if action.ID == "forward_enquiry" {
				if err := helpers.ShowForwardView(payload.TriggerID, enquiry); err != nil {
					panic(err)
				}
			} else if action.ID == "trash_enquiry" {
				if err := helpers.TrashEnquiry(payload.User.ID, enquiry); err != nil {
					panic(err)
				}
			}
		}
	} else if payload.Type == "view_submission" && payload.View.PrivateMetadata != "" {
		id, _ := strconv.Atoi(payload.View.PrivateMetadata)
		var enquiry models.Enquiry
		services.GormDb().Find(&enquiry, id)

		if payload.View.CallbackID == "follow_up_enquiry" {
			if err := helpers.RecordFollowUp(payload.User.ID, enquiry, payload.View.State.Values); err != nil {
				panic(err)
			}
		} else if payload.View.CallbackID == "follow_up_forward" {
			if err := helpers.ForwardEnquiry(payload.User.ID, enquiry, payload.View.State.Values); err != nil {
				panic(err)
			}
		}
	}

	c.Status(http.StatusOK)
}

func SlackOptions(c *gin.Context) {
	if !services.VerifyRequestFromSlack(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Signature from Slack is not valid."})
		return
	}

	var input models.SlackPayload
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload models.SlackOptionsPayload
	if err := json.Unmarshal([]byte(input.Payload), &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Type == "block_suggestion" && payload.ActionID == "follow_up_enquiry_member_input" {
		options, err := helpers.LoadForwardToRecipientOptions(payload.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"options": options,
		})
	}

	c.Status(http.StatusOK)
}
