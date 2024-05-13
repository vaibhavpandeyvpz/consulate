package models

type ExotelCall struct {
	CallSid              string `json:"CallSid" required:"true"`
	EventType            string `json:"EventType" required:"true"`
	DateCreated          string `json:"DateCreated" required:"true"`
	DateUpdated          string `json:"DateUpdated" required:"true"`
	Status               string `json:"Status" required:"true"`
	From                 string `json:"From" required:"true"`
	To                   string `json:"To" required:"true"`
	PhoneNumberSid       string `json:"PhoneNumberSid" required:"true"`
	StartTime            string `json:"StartTime" required:"true"`
	EndTime              string `json:"EndTime" required:"true"`
	Direction            string `json:"Direction" required:"true"`
	RecordingUrl         string `json:"RecordingUrl"`
	ConversationDuration int    `json:"ConversationDuration" required:"true"`
	CustomField          string `json:"CustomField"`
}

type ExotelCustomField struct {
	UserID    string `json:"user_id"`
	EnquiryID uint   `json:"enquiry_id"`
}
