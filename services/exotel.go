package services

import (
	"consulate/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func CallWithExotel(user string, enquiry models.Enquiry, from string, to string) (err error) {
	endpoint := fmt.Sprintf(
		"https://%s:%s@%s/v1/Accounts/%s/Calls/connect",
		GetConfig().Exotel.ApiKey,
		GetConfig().Exotel.ApiToken,
		GetConfig().Exotel.Domain,
		GetConfig().Exotel.AccountSid,
	)

	form := url.Values{}
	form.Add("From", from)
	form.Add("To", to)
	form.Add("CallerId", GetConfig().Exotel.CallerId)
	form.Add("CallType", "trans")
	form.Add("Record", "true")
	customField := models.ExotelCustomField{EnquiryID: enquiry.ID, UserID: user}
	bytes, err := json.Marshal(customField)
	if err != nil {
		return err
	}

	form.Add("CustomField", string(bytes))
	form.Add("StatusCallback", GetConfig().Exotel.StatusCallbackUrl)
	form.Add("StatusCallbackEvents[0]", "terminal")
	_, err = http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	return
}
