package kecpws

import (
	"time"

	kchan "github.com/fourdim/kecp/modules/kecp-channel"
	kecpcrypto "github.com/fourdim/kecp/modules/kecp-crypto"
)

const (
	roomLiveCheckWait = 30 * time.Second
)

type Room struct {
	RoomID string

	// Creator's key
	ManagementKey string

	// Registry
	registry *Registry

	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast *kchan.Channel[[]byte]

	// Register requests from the clients.
	register *kchan.Channel[*Client]

	// Unregister requests from clients.
	unregister *kchan.Channel[*Client]

	// The status returned after register.
	created chan bool

	// Channel for self destruction.
	selfDestruction chan bool
}

func (reg *Registry) NewRoom(managementKey string) (string, error) {
	roomID := kecpcrypto.GenerateToken()
	room := &Room{
		RoomID:          roomID,
		ManagementKey:   managementKey,
		registry:        reg,
		broadcast:       kchan.New[[]byte](),
		register:        kchan.New[*Client](),
		unregister:      kchan.New[*Client](),
		clients:         make(map[string]*Client),
		created:         make(chan bool),
		selfDestruction: make(chan bool),
	}
	room.registry.register.Write(room)
	checker := time.NewTimer(clientJoinedCheckWait)
	defer checker.Stop()
	select {
	case <-room.created:
	case <-checker.C:
		return "", ErrCanNotCreateTheRoom
	}
	go room.run()
	return roomID, nil
}

func (room *Room) run() {
	// Only this goroutine can access
	// Room.clients
	checker := time.NewTimer(roomLiveCheckWait)
	defer func() {
		checker.Stop()
		room.registry.unregister.Write(room)
	}()
	for {
		select {
		case client := <-room.register.Read():
			if previousClient, ok := room.clients[client.ClientKey]; ok {
				previousClient.selfDestruction <- true
			}
			room.clients[client.ClientKey] = client
			client.joined <- true
		case clientUnregistered := <-room.unregister.Read():
			if client, ok := room.clients[clientUnregistered.ClientKey]; ok {
				if client == clientUnregistered {
					delete(room.clients, clientUnregistered.ClientKey)
				}
				close(clientUnregistered.send)
				close(clientUnregistered.joined)
				close(clientUnregistered.selfDestruction)
			}
			if len(room.clients) == 0 {
				return
			}
		case message := <-room.broadcast.Read():
			for clientKey, client := range room.clients {
				select {
				case client.send <- message:
				default:
					delete(room.clients, clientKey)
					close(client.send)
					close(client.joined)
					close(client.selfDestruction)
				}
			}
			if len(room.clients) == 0 {
				return
			}
		// Delete the room if no one joins.
		case <-checker.C:
			if len(room.clients) == 0 {
				return
			}
		case <-room.selfDestruction:
			for _, client := range room.clients {
				client.selfDestruction <- true
			}
		}
	}
}
