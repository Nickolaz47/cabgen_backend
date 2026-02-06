package events

type UserRegisteredPayload struct {
	RegisteredUsername string `json:"username"`
}