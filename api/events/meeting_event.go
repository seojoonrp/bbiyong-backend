// api/events/meeting_event.go

package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/api/ws"
	"github.com/seojoonrp/bbiyong-backend/models"
)

func StartMeetingWorker(eventChan <-chan models.MeetingEvent, chatService services.ChatService, hub *ws.Hub) {
	for event := range eventChan {
		go func(e models.MeetingEvent) {
			ctx := context.Background()

			msg, err := chatService.SaveSystemMessage(ctx, e.MeetingID, e.UserID, e.Type)
			if err != nil {
				log.Printf("Failed to save system message: %v", err)
				return
			}

			payload, _ := json.Marshal(msg)
			hub.Broadcast <- ws.MessagePayload{
				MeetingID: e.MeetingID,
				Data:      payload,
			}
		}(event)
	}
}
