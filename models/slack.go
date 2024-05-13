package models

type SlackAction struct {
	ID             string `json:"action_id"`
	SelectedOption struct {
		Value string `json:"value"`
	} `json:"selected_option"`
	Value string `json:"value"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackInteractionPayload struct {
	Type      string        `json:"type"`
	TriggerID string        `json:"trigger_id"`
	Team      SlackTeam     `json:"team"`
	User      SlackUser     `json:"user"`
	Channel   SlackChannel  `json:"channel"`
	Message   SlackMessage  `json:"message"`
	Actions   []SlackAction `json:"actions"`
	View      SlackView     `json:"view"`
}

type SlackPayload struct {
	Payload string `form:"payload"`
}

type SlackMessage struct {
	Team string `json:"team"`
	User string `json:"user"`
	Ts   string `json:"ts"`
}

type SlackOptionsPayload struct {
	Type     string    `json:"type"`
	Team     SlackTeam `json:"team"`
	User     SlackUser `json:"user"`
	BlockID  string    `json:"block_id"`
	ActionID string    `json:"action_id"`
	Value    string    `json:"value"`
	View     SlackView `json:"view"`
}

type SlackTeam struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type SlackUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackView struct {
	CallbackID      string `json:"callback_id"`
	PrivateMetadata string `json:"private_metadata"`
	State           struct {
		Values SlackViewStateValues `json:"values"`
	} `json:"state"`
}

type SlackViewStateValues map[string]map[string]map[string]any
