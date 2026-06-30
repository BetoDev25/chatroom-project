package main

type Message struct {
	Type           string `json:"type`      //"join", "leave", "message"
	Room           string `json:"room"`     //target room name
	Username       string `json:"username"` //display name
	Content        string `json:"content"`  //message body
	SkipBroadcast  bool   `json:"skipBroadcast,omitempty"`
	ConversationID string `json:"conversation_id"`
}
