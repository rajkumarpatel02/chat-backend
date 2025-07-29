// ✅ models/message.go
package models

type Message struct {
	Type    string `json:"type"`     // "private" or "public"
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
