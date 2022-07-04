package kecpsignal

import kchan "github.com/fourdim/kecp/modules/kecp-channel"

type Registry struct {
	rooms map[string]*Room

	// register is written by the rooms
	register *kchan.Channel[*Room]

	// unregister is written by the rooms
	unregister *kchan.Channel[*Room]

	// roomQuery is written by the clients
	roomQuery *kchan.Channel[*roomQuery]

	// roomDeletionRequest
	roomDeletionRequest chan *roomDeletion
}

func NewRegistry() *Registry {
	reg := &Registry{
		rooms:               make(map[string]*Room),
		register:            kchan.New[*Room](),
		unregister:          kchan.New[*Room](),
		roomQuery:           kchan.New[*roomQuery](),
		roomDeletionRequest: make(chan *roomDeletion),
	}
	go reg.run()
	return reg
}

func (reg *Registry) run() {
	// Only this goroutine can access
	// Registry.rooms
	for {
		select {
		case room := <-reg.register.Read():
			reg.rooms[room.RoomID] = room
			room.created <- true
		case room := <-reg.unregister.Read():
			if _, ok := reg.rooms[room.RoomID]; ok {
				delete(reg.rooms, room.RoomID)
				room.broadcast.Close()
				room.register.Close()
				room.unregister.Close()
				close(room.created)
				close(room.selfDestruction)
			}
		case roomQuery := <-reg.roomQuery.Read():
			if room, ok := reg.rooms[roomQuery.roomID]; ok {
				roomQuery.room <- room
			} else {
				roomQuery.room <- nil
			}
			close(roomQuery.room)
		case roomDele := <-reg.roomDeletionRequest:
			if room, ok := reg.rooms[roomDele.roomID]; ok && room.MgtKey == roomDele.mgtKey {
				room.selfDestruction <- true
			}
		}
	}
}

type roomQuery struct {
	roomID string
	room   chan *Room
}

func (reg *Registry) GetRoom(roomID string) *Room {
	roomQuery := &roomQuery{
		roomID: roomID,
		room:   make(chan *Room),
	}
	reg.roomQuery.Write(roomQuery)
	room := <-roomQuery.room
	return room
}

type roomDeletion struct {
	roomID string
	mgtKey string
}

func (reg *Registry) DeleteRoom(roomID string, managementKey string) {
	reg.roomDeletionRequest <- &roomDeletion{roomID: roomID, mgtKey: managementKey}
}
