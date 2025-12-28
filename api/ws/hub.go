// api/ws/hub.go

package ws

import "sync"

type Hub struct {
	Rooms      map[string]map[*Client]bool // 어떤 미팅ID에 어떤 클라이언트들이 연결되어 있는지
	Broadcast  chan MessagePayload         // 메시지 전달 채널
	Register   chan *Client                // 나 여기 들어간다
	Unregister chan *Client                // 나 나간다
	mu         sync.Mutex
}

type MessagePayload struct {
	MeetingID string
	Data      []byte // JSON-encoded
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan MessagePayload),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// 무한루프 돌면서 브로드캐스팅 처리
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Rooms[client.MeetingID] == nil {
				h.Rooms[client.MeetingID] = make(map[*Client]bool)
			}
			h.Rooms[client.MeetingID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.MeetingID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Rooms, client.MeetingID)
					}
				}
			}
			h.mu.Unlock()

		case payload := <-h.Broadcast:
			h.mu.Lock()
			clients := h.Rooms[payload.MeetingID]
			for client := range clients {
				select {
				case client.Send <- payload.Data:
				default: // 인터넷 연결 불안정 등으로 전송이 안될 때 -> 넌 나가라!
					close(client.Send)
					delete(clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
