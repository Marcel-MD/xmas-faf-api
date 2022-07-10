package websockets

import (
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/rs/zerolog/log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan models.Message
	register   chan subscription
	unregister chan subscription
	service    services.IMessageService
}

var once sync.Once
var h hub

func InitHub() {
	once.Do(func() {
		log.Info().Msg("Initializing websocket hub")
		h = hub{
			rooms:      make(map[string]map[*connection]bool),
			broadcast:  make(chan models.Message),
			register:   make(chan subscription),
			unregister: make(chan subscription),
			service:    services.GetMessageService(),
		}
		go h.run()
	})
}

func (h *hub) run() {
	log.Info().Msg("Starting websocket hub")
	for {
		select {

		case s := <-h.register:
			log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Register user from room connection")
			connections := h.rooms[s.roomID]
			if connections == nil {
				log.Debug().Str("room_id", s.roomID).Msg("Creating room connection")
				connections = make(map[*connection]bool)
				h.rooms[s.roomID] = connections
			}
			h.rooms[s.roomID][s.conn] = true

		case s := <-h.unregister:
			connections := h.rooms[s.roomID]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Unregister user from room connection")
					close(s.conn.send)
					delete(connections, s.conn)
					if len(connections) == 0 {
						log.Debug().Str("room_id", s.roomID).Msg("Deleting room connection")
						delete(h.rooms, s.roomID)
					}
				}
			}

		case m := <-h.broadcast:
			log.Debug().Str("room_id", m.RoomID).Str("msg_id", m.ID).Msg("Broadcasting message")
			connections := h.rooms[m.RoomID]
			for c := range connections {
				select {
				case c.send <- m:
				default:
					log.Warn().Str("user_id", m.UserID).Str("room_id", m.RoomID).Msg("Closing user connection")
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						log.Debug().Str("room_id", m.RoomID).Msg("Deleting room connection")
						delete(h.rooms, m.RoomID)
					}
				}
			}
		}
	}
}
