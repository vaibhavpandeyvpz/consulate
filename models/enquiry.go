package models

import (
	"gorm.io/gorm"
)

type Enquiry struct {
	gorm.Model
	Name           string `form:"name"  json:"name"`
	Email          string `form:"email" json:"email"`
	Phone          string `form:"phone" json:"phone" binding:"required"`
	Message        string `form:"message" json:"message"`
	Status         string
	SlackMessageTs string
	FollowUps      []FollowUp
}

type FollowUp struct {
	gorm.Model
	Notes     string
	EnquiryID int
	Enquiry   Enquiry
}
