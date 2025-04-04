// server/reply.go
package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mizzobg/mailpit-send/internal/smtpclient"
)

// ReplyRequest represents the JSON payload for a reply.
type ReplyRequest struct {
	MessageID string `json:"message_id"`        // optional: original message ID
	Subject   string `json:"subject,omitempty"` // if empty, a default will be used
	Body      string `json:"body"`              // reply body
	// You could also add fields like Recipient if desired.
}

// ReplyHandler handles POST requests to send a reply.
func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	var req ReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// (Optional) Retrieve the original message by req.MessageID if needed.
	// For example:
	// originalMsg, err := messageStore.GetMessage(req.MessageID)
	// if err != nil { … }

	// Create a new reply message.
	replyMsg := smtpclient.NewMessage()
	// Set a default sender (could be a system “noreply” address)
	replyMsg.From = "noreply@example.com"

	// Set the subject – prepend "Re:" if not already present.
	if req.Subject != "" {
		replyMsg.Subject = req.Subject
	} else {
		replyMsg.Subject = "Re: (no subject)"
	}
	replyMsg.Body = req.Body

	// For this example, we hardcode the recipient.
	// In practice, you might extract the original sender from the looked-up message.
	replyMsg.To = []string{"original.sender@example.com"}

	// Optionally, add a header to indicate this is a reply.
	// If you had the original message ID, you might do:
	// replyMsg.Headers["In-Reply-To"] = originalMsg.MessageID

	if err := smtpclient.Send(replyMsg); err != nil {
		log.Printf("Error sending reply: %v", err)
		http.Error(w, "Failed to send reply", http.StatusInternalServerError)
		return
	}

	// Return a JSON response indicating success.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}
