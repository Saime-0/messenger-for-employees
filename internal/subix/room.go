package subix

type Chats map[ID]*Chat

type Chat struct {
	ID      int
	members Rooms
}

type Rooms map[ID]*Room

type Room struct {
	RoomID            int
	Empls             Employees
	clientsWithEvents ClientsWithEvents
}

func (s *Subix) CreateRoomIfNotExists(roomID, chatID int) *Room {
	room, ok := s.rooms[roomID]
	if !ok {
		room = &Room{
			RoomID:            roomID,
			Empls:             Employees{},
			clientsWithEvents: make(ClientsWithEvents),
		}

		s.rooms[roomID] = room
	}
	return room
}

func (s *Subix) DeleteRoom(roomID int) {
	room, ok := s.rooms[roomID]
	if ok { // если вдруг не удается найти, то просто пропускаем
		delete(s.rooms, roomID)      // удаление из глобальной мапы
		room.clientsWithEvents = nil // на всякий случай заnullяем мапу

		emp, ok := s.employees[room.RoomID]
		if ok {
			delete(emp.rooms, room.RoomID)
		}
	}
}
