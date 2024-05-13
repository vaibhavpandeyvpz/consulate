package models

type SlackAction struct {
	ID    string `json:"action_id" required:"true"`
	Value string `json:"value" required:"true"`
}

type SlackChannel struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name" required:"true"`
}

type SlackInteractionPayload struct {
	Type      string        `json:"type" required:"true"`
	TriggerID string        `json:"trigger_id" required:"true"`
	Team      SlackTeam     `json:"team" required:"true"`
	User      SlackUser     `json:"user" required:"true"`
	Channel   SlackChannel  `json:"channel"`
	Message   SlackMessage  `json:"message"`
	Actions   []SlackAction `json:"actions"`
	View      SlackView     `json:"view"`
}

type SlackPayload struct {
	Payload string `form:"payload" required:"true"`
}

type SlackMessage struct {
	Team string `json:"team" required:"true"`
	User string `json:"user" required:"true"`
	Ts   string `json:"ts" required:"true"`
}

type SlackOptionsPayload struct {
	Type     string    `json:"type" required:"true"`
	Team     SlackTeam `json:"team" required:"true"`
	User     SlackUser `json:"user" required:"true"`
	BlockID  string    `json:"block_id" required:"true"`
	ActionID string    `json:"action_id" required:"true"`
	Value    string    `json:"value" required:"true"`
	View     SlackView `json:"view"`
}

type SlackTeam struct {
	ID     string `json:"id" required:"true"`
	Domain string `json:"domain" required:"true"`
}

type SlackUser struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name" required:"true"`
}

type SlackView struct {
	CallbackID      string `json:"callback_id"`
	PrivateMetadata string `json:"private_metadata"`
	State           struct {
		Values SlackViewStateValues `json:"values"`
	} `json:"state"`
}

type SlackViewStateValues map[string]map[string]map[string]any
