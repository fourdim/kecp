package kecpsignal

import (
	"time"

	kchan "github.com/fourdim/kecp/modules/kecp-channel"
	kecpcrypto "github.com/fourdim/kecp/modules/kecp-crypto"
	kecpmsg "github.com/fourdim/kecp/modules/kecp-msg"
	kecpvalidate "github.com/fourdim/kecp/modules/kecp-validate"
)

const (
	roomLiveCheckWait = 30 * time.Second
)

type Room struct {
	RoomID string

	// Creator's key
	MgtKey string

	// Registry
	registry *Registry

	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast *kchan.Channel[*kecpmsg.Message]

	// Inbound messages from the clients.
	forward *kchan.Channel[*kecpmsg.Message]

	// Register requests from the clients.
	register *kchan.Channel[*Client]

	// Unregister requests from clients.
	unregister *kchan.Channel[*Client]

	// The status returned after register.
	created chan bool

	// Channel for self destruction.
	selfDestruction chan bool
}

func (reg *Registry) NewRoom(managementKey string) string {
	if !kecpvalidate.IsAValidCryptoKey(managementKey) {
		return ""
	}
	roomID := kecpcrypto.GenerateRoomID()
	room := &Room{
		RoomID:          roomID,
		MgtKey:          managementKey,
		registry:        reg,
		broadcast:       kchan.New[*kecpmsg.Message](),
		forward:         kchan.New[*kecpmsg.Message](),
		register:        kchan.New[*Client](),
		unregister:      kchan.New[*Client](),
		clients:         make(map[string]*Client),
		created:         make(chan bool),
		selfDestruction: make(chan bool),
	}
	room.registry.register.Write(room)
	select {
	case <-room.created:
	}
	go room.run()
	return roomID
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
			var joined = true
			for _, eachClient := range room.clients {
				// Same name, but not the same client.
				if client.name == eachClient.name && client.clientKey != eachClient.clientKey {
					joined = false
					break
				}
			}
			if !joined {
				client.joined <- false
				break
			}
			if previousClient, ok := room.clients[client.clientKey]; ok {
				previousClient.selfDestruction <- true
				broadcast(room, kecpmsg.NewLeaveMsg(previousClient.name, previousClient.clientKey))
			}
			room.clients[client.clientKey] = client
			client.joined <- true
			var names []string
			for _, eachClient := range room.clients {
				names = append(names, eachClient.name)
			}
			client.send <- kecpmsg.NewListMsg(names)
			broadcast(room, kecpmsg.NewJoinMsg(client.name, client.clientKey))
		case clientUnregistered := <-room.unregister.Read():
			var replace bool
			if client, ok := room.clients[clientUnregistered.clientKey]; ok {
				if client == clientUnregistered {
					replace = false
					delete(room.clients, clientUnregistered.clientKey)
				} else {
					replace = true
				}
			}
			close(clientUnregistered.send)
			close(clientUnregistered.joined)
			close(clientUnregistered.selfDestruction)
			if !replace {
				broadcast(room, kecpmsg.NewLeaveMsg(clientUnregistered.name, clientUnregistered.clientKey))
			}
			if len(room.clients) == 0 {
				return
			}
		case message := <-room.forward.Read():
			forward(room, message)
			if len(room.clients) == 0 {
				return
			}
		case message := <-room.broadcast.Read():
			broadcast(room, message)
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

func broadcast(room *Room, message *kecpmsg.Message) {
	for clientKey, client := range room.clients {
		if message.ExceptClientKey == clientKey {
			continue
		}
		sendToSingleClient(room, client, message)
	}
}

func forward(room *Room, message *kecpmsg.Message) {
	for _, client := range room.clients {
		if message.Target == client.name {
			sendToSingleClient(room, client, message)
			break
		}
	}
}

func sendToSingleClient(room *Room, client *Client, message *kecpmsg.Message) {
	select {
	case client.send <- message:
	default:
		delete(room.clients, client.clientKey)
		close(client.send)
		close(client.joined)
		close(client.selfDestruction)
	}
}
