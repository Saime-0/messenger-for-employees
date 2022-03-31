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

func (s *Subix) CreateRoomIfNotExists(roomID int) *Room {
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
		delete(s.rooms, roomID) // удаление из глобальной мапы
		for _, emp := range room.Empls {
			delete(emp.rooms, roomID)
			//delete(s.employees[emp.Client.EmployeeID].rooms, roomID)
		}
		room.clientsWithEvents = nil // на всякий случай заnullяем мапу
	}
}

func (s *Subix) DeleteMember(roomID int, empID int) {
	emp, ok := s.employees[empID]
	if ok {
		room, ok := s.rooms[roomID]
		if ok {
			for _, client := range emp.clients { // удалить клиентов из комнаты на основе клиентов сотрудника
				delete(room.clientsWithEvents, client.sessionKey)
			}
			delete(room.Empls, empID) // удалить сотрудника из комнаты

			if len(room.Empls) == 0 {
				s.DeleteRoom(roomID)
			}
		}
		delete(emp.rooms, roomID) // удалить комнату из сотрудника

	}
}
