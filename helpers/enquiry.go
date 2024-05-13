package helpers

import (
	"consulate/models"
	"consulate/services"
	"encoding/json"
	"fmt"
	"github.com/ggwhite/go-masker"
	"github.com/samber/lo"
	"github.com/slack-go/slack"
	"slices"
	"strconv"
	"strings"
)

func ForwardEnquiry(user string, enquiry models.Enquiry, values models.SlackViewStateValues) (err error) {
	sc := services.SlackClient()
	_ = sc.AddReaction("dart", slack.ItemRef{
		Channel:   services.GetConfig().Slack.Channels.Enquiries,
		Timestamp: enquiry.SlackMessageTs,
	})

	enquiry.Status = "conversation"
	services.GormDb().Save(&enquiry)

	selectedOption := values["follow_up_forward_member"]["follow_up_enquiry_member_input"]["selected_option"].(map[string]any)
	recipient := selectedOption["value"]
	_, _, err = sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionText(
			fmt.Sprintf("<@%s> has *forwarded* this enquiry to <@%s>", user, recipient), false,
		),
		slack.MsgOptionTS(enquiry.SlackMessageTs),
	)
	if err != nil {
		return
	}

	permalink, err := sc.GetPermalink(&slack.PermalinkParameters{
		Channel: services.GetConfig().Slack.Channels.Enquiries,
		Ts:      enquiry.SlackMessageTs,
	})
	if err != nil {
		return
	}

	var fields []*slack.TextBlockObject
	if len(enquiry.Email) > 0 {
		fields = append(fields, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Email address*\n<mailto:%s|%s>", enquiry.Email, enquiry.Email),
		})
	}

	fields = append(fields, &slack.TextBlockObject{
		Type: "mrkdwn",
		Text: fmt.Sprintf("*Phone number*\n<tel:%s|%s>", enquiry.Phone, enquiry.Phone),
	})

	_, _, err = sc.PostMessage(
		recipient.(string),
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					"mrkdwn",
					fmt.Sprintf("<@%s> has forwarded you <%s|this enquiry>. Please find the direct contact details below:", user, permalink),
					false,
					false,
				),
				fields,
				nil,
			),
		),
		slack.MsgOptionDisableLinkUnfurl(),
		slack.MsgOptionDisableMediaUnfurl(),
	)

	return
}

func HandleNewEnquiry(enquiry models.Enquiry, ip string) (err error) {
	services.GormDb().Save(&enquiry)

	var fields []*slack.TextBlockObject
	if len(enquiry.Name) > 0 {
		fields = append(fields, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: "*Name*\n" + enquiry.Name,
		})
	}

	if len(enquiry.Email) > 0 {
		fields = append(fields, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: "*Email address*\n" + masker.Email(enquiry.Email),
		})
	}

	fields = append(fields, &slack.TextBlockObject{
		Type: "mrkdwn",
		Text: "*Phone number*\n" + masker.Mobile(enquiry.Phone),
	})

	ipd, err := services.FindIpLocation(ip)
	if err != nil {
		return
	}

	fields = append(fields, &slack.TextBlockObject{
		Type: "mrkdwn",
		Text: "*Location*\n" + fmt.Sprintf("%s, %s", ipd.Location.City, ipd.Location.Country),
	})

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject(
				"plain_text",
				"Someone just :inbox_tray: submitted an enquiry.",
				false,
				false,
			),
			fields,
			nil,
		),
	}

	msg := strings.TrimSpace(enquiry.Message)
	if len(msg) > 0 {
		blocks = append(blocks, slack.NewSectionBlock(
			slack.NewTextBlockObject(
				"mrkdwn",
				fmt.Sprintf("*Message:*\n%s", msg),
				false,
				false,
			),
			nil,
			nil,
		))
	}

	blocks = append(blocks, slack.NewActionBlock(
		"",
		slack.NewButtonBlockElement(
			"call_now",
			strconv.Itoa(int(enquiry.ID)),
			slack.NewTextBlockObject(
				"plain_text",
				"Call now",
				false,
				false,
			),
		),
		slack.NewButtonBlockElement(
			"follow_up_enquiry",
			strconv.Itoa(int(enquiry.ID)),
			slack.NewTextBlockObject(
				"plain_text",
				"Follow up",
				false,
				false,
			),
		),
		slack.ButtonBlockElement{
			Type:     slack.METButton,
			ActionID: "forward_enquiry",
			Text: slack.NewTextBlockObject(
				"plain_text",
				"Forward",
				false,
				false,
			),
			Style: "primary",
			Value: strconv.Itoa(int(enquiry.ID)),
		},
		slack.ButtonBlockElement{
			Type:     slack.METButton,
			ActionID: "trash_enquiry",
			Text: slack.NewTextBlockObject(
				"plain_text",
				"Trash",
				false,
				false,
			),
			Confirm: slack.NewConfirmationBlockObject(
				slack.NewTextBlockObject("plain_text", "Really?", false, false),
				slack.NewTextBlockObject("plain_text", "This action will mark the enquiry as trashed.", false, false),
				slack.NewTextBlockObject("plain_text", "Yes", false, false),
				slack.NewTextBlockObject("plain_text", "Cancel", false, false),
			),
			Style: "danger",
			Value: strconv.Itoa(int(enquiry.ID)),
		},
	))

	sc := services.SlackClient()
	_, ts, err := sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionBlocks(blocks...),
	)

	enquiry.SlackMessageTs = ts

	services.GormDb().Save(&enquiry)

	return
}

func LoadForwardToRecipientOptions(q string) (options []slack.OptionBlockObject, err error) {
	sc := services.SlackClient()
	users, err := sc.GetUsers()
	if err != nil {
		return
	}

	allowedUsers := lo.Filter(users, func(user slack.User, index int) bool {
		return slices.Contains(services.GetConfig().Recipients, user.ID)
	})

	matchingUsers := lo.Filter(allowedUsers, func(user slack.User, index int) bool {
		return strings.Contains(strings.ToLower(user.RealName), strings.ToLower(q))
	})

	options = lo.Map(matchingUsers, func(user slack.User, index int) slack.OptionBlockObject {
		return slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", user.RealName, false, false),
			Value: user.ID,
		}
	})

	return
}

func PlacePhoneCall(user string, enquiry models.Enquiry) (err error) {
	sc := services.SlackClient()
	info, err := sc.GetUserInfo(user)
	if err != nil {
		return
	}

	if info.Profile.Phone == "" {
		_, err = sc.PostEphemeral(
			services.GetConfig().Slack.Channels.Enquiries,
			user,
			slack.MsgOptionText("Please add a phone number to your Slack profile to use this feature.", false),
		)

		return
	}

	_ = sc.AddReaction("call_me_hand::skin-tone-4", slack.ItemRef{
		Channel:   services.GetConfig().Slack.Channels.Enquiries,
		Timestamp: enquiry.SlackMessageTs,
	})

	enquiry.Status = "conversation"
	services.GormDb().Save(&enquiry)

	_, _, err = sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionText(
			fmt.Sprintf("<@%s> has *called* this enquiry.", user), false,
		),
		slack.MsgOptionTS(enquiry.SlackMessageTs),
	)
	if err != nil {
		return
	}

	err = services.CallWithExotel(user, enquiry, info.Profile.Phone, enquiry.Phone)

	return
}

func RecordFollowUp(user string, enquiry models.Enquiry, values models.SlackViewStateValues) (err error) {
	sc := services.SlackClient()
	_ = sc.AddReaction("parrot", slack.ItemRef{
		Channel:   services.GetConfig().Slack.Channels.Enquiries,
		Timestamp: enquiry.SlackMessageTs,
	})

	enquiry.Status = "conversation"
	services.GormDb().Save(&enquiry)

	notes := values["follow_up_enquiry_notes"]["follow_up_enquiry_notes_input"]["value"]
	followUp := &models.FollowUp{Notes: notes.(string)}
	err = services.GormDb().Model(&enquiry).Association("FollowUps").Append(followUp)
	if err != nil {
		return
	}

	_, _, err = sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionText(
			fmt.Sprintf("<@%s> has *recorded a follow up* on this enquiry i.e., %s", user, notes), false,
		),
		slack.MsgOptionTS(enquiry.SlackMessageTs),
	)

	return
}

func ShowFollowUpView(triggerId string, enquiry models.Enquiry) (err error) {
	sc := services.SlackClient()
	_, err = sc.OpenView(triggerId, slack.ModalViewRequest{
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					BlockID: "follow_up_enquiry_notes",
					Element: slack.PlainTextInputBlockElement{
						ActionID:    "follow_up_enquiry_notes_input",
						Multiline:   true,
						Placeholder: slack.NewTextBlockObject("plain_text", "Details of conversation…", false, false),
						Type:        "plain_text_input",
					},
					Label: slack.NewTextBlockObject("plain_text", "Notes", false, false),
					Type:  "input",
				},
			},
		},
		CallbackID:      "follow_up_enquiry",
		PrivateMetadata: strconv.Itoa(int(enquiry.ID)),
		Submit:          slack.NewTextBlockObject("plain_text", "Save", false, false),
		Title:           slack.NewTextBlockObject("plain_text", "Record follow up", false, false),
		Type:            "modal",
	})

	return
}

func ShowForwardView(triggerId string, enquiry models.Enquiry) (err error) {
	sc := services.SlackClient()
	minQuery := 1
	_, err = sc.OpenView(triggerId, slack.ModalViewRequest{
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					BlockID: "follow_up_forward_member",
					Element: slack.SelectBlockElement{
						ActionID:       "follow_up_enquiry_member_input",
						MinQueryLength: &minQuery,
						Placeholder:    slack.NewTextBlockObject("plain_text", "Search for a user…", false, false),
						Type:           "external_select",
					},
					Label: slack.NewTextBlockObject("plain_text", "Member", false, false),
					Type:  "input",
				},
			},
		},
		CallbackID:      "follow_up_forward",
		PrivateMetadata: strconv.Itoa(int(enquiry.ID)),
		Submit:          slack.NewTextBlockObject("plain_text", "Send", false, false),
		Title:           slack.NewTextBlockObject("plain_text", "Forward to", false, false),
		Type:            "modal",
	})

	return
}

func SendCallRecording(call models.ExotelCall) (err error) {
	var customField models.ExotelCustomField
	if err = json.Unmarshal([]byte(call.CustomField), &customField); err != nil {
		return
	}

	var enquiry models.Enquiry
	services.GormDb().Find(&enquiry, customField.EnquiryID)

	sc := services.SlackClient()
	_, _, err = sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					"mrkdwn",
					fmt.Sprintf(
						"<@%s> has finished a %ds *phone call* with this enquiry.",
						customField.UserID,
						call.ConversationDuration,
					),
					false,
					false,
				),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"",
				slack.NewButtonBlockElement(
					"",
					"",
					slack.NewTextBlockObject(
						"plain_text",
						"Download",
						false,
						false,
					),
				).WithURL(call.RecordingUrl),
			),
		),
		slack.MsgOptionTS(enquiry.SlackMessageTs),
	)

	return
}

func TrashEnquiry(user string, enquiry models.Enquiry) (err error) {
	sc := services.SlackClient()
	_ = sc.AddReaction("x", slack.ItemRef{
		Channel:   services.GetConfig().Slack.Channels.Enquiries,
		Timestamp: enquiry.SlackMessageTs,
	})

	enquiry.Status = "trashed"
	services.GormDb().Save(&enquiry)

	_, _, err = sc.PostMessage(
		services.GetConfig().Slack.Channels.Enquiries,
		slack.MsgOptionText(
			fmt.Sprintf("<@%s> has *trashed* this enquiry.", user), false,
		),
		slack.MsgOptionTS(enquiry.SlackMessageTs),
	)

	return
}
