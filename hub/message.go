package cameraHub

import (
	"encoding/json"
)
// Message is json sent from client
type Message struct {
	Type       string      `json:"type"` // Clients, AddName, Offer, Answer
	Data       interface{} `json:"data"`
}

// func (m *Message) Json() ([]byte) {
//   j, _ := json.Marshal(m)
//   return j
// }

func unmarshalMessage(bytes []byte) (*Message) {
  var message Message
  json.Unmarshal(bytes, &message)
  return &message
}
