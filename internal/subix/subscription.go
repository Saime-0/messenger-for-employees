package subix

import (
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
)

func (s *Subix) Sub(userID int, sessionKey Key, expAt int64) (*Client, error) {
	_, ok := s.clients[sessionKey]
	if ok { // если ключ существует, то по-хорошему клиент должен повторить соединение с другим ключом
		return nil, cerrors.New("sessionKey already in use, it is not possible to create a new connection")
	}

	// websocket connection = сессия = клиент
	client := &Client{
		EmployeeID:       userID,
		Ch:               make(chan *model.SubscriptionBody),
		sessionExpiresAt: expAt,
		sessionKey:       sessionKey,
	}
	s.clients[sessionKey] = client

	// планируем пометку и дальнейшее удаление клиента, если его токен истечет
	err := s.scheduleMarkClient(client, expAt)
	if err != nil {
		delete(s.clients, sessionKey)
		return nil, cerrors.New("не удалось создать сессию")
	}

	// user
	user := s.CreateEmployeeIfNotExists(userID)
	user.clients[sessionKey] = client

	return client, nil
}

var allUsefulEventTypes = model.AllEventType[1:len(model.AllEventType)]

func (s *Subix) ModifyCollection(employeeID int, sessionKey Key, roomIDs []ID, action model.EventSubjectAction, listenEvents []model.EventType) error {
	client, ok := s.clients[sessionKey]
	if !ok { // если ключа не существует, то клиент должен подписаться
		return cerrors.New("no session with this key was found")
	}
	emp, ok := s.employees[client.EmployeeID]
	if !ok { // сотрудник создается во время подписки, либо он случайно удалился, либо
		return cerrors.New("client is not associated with any employee")
	}
	for _, event := range listenEvents {
		if event == model.EventTypeAll {
			listenEvents = allUsefulEventTypes
			break
		}
	}
	if action == model.EventSubjectActionAdd { // если мемберс добавляет пачку событий
		for _, roomID := range roomIDs {
			room := s.CreateRoomIfNotExists(roomID)                    // достаем комнату из активных(те на которые подписаны клиенты) нужную комнату
			clientWithEvents, ok := room.clientsWithEvents[sessionKey] // ищем сессию нужного клиента в комнате
			if !ok {                                                   // если клиент еще не прослушивает этого участника, то заставляем слушать
				clientWithEvents = &ClientWithEvents{
					Client: client,
					Events: make(EventCollection),
				}
				room.clientsWithEvents[sessionKey] = clientWithEvents
			}
			for _, event := range listenEvents {
				clientWithEvents.Events[event] = true // добавил тип ивента который теперь будет отправляться клиенту(прослушиваться им)
			}

			// add room to emp rooms
			emp.rooms[roomID] = room // даже если у пользователя уже есть комната с таким id то все равно добавляем(не даст никакого эффекта)
			// and emp to room emps
			room.Empls[emp.EmpID] = emp
		}

	} else if action == model.EventSubjectActionDelete {
		for _, roomID := range roomIDs {

			room, ok := s.rooms[roomID]
			if !ok {
				continue // пропускаем если комнаты не существует (ее никто не прослушивает)
			} else {
				clientWithEvents, ok := room.clientsWithEvents[sessionKey]
				if !ok {
					break // клиент не слушает комнату, а значит и удалять ничего не нужно
				}
				for _, event := range listenEvents {
					delete(clientWithEvents.Events, event)
				}

				if len(clientWithEvents.Events) == 0 { // если удалятся все события которые прослушивал клиент, то ..
					delete(room.clientsWithEvents, sessionKey) // .. Удаляем клиента из комнаты
					if len(room.clientsWithEvents) == 0 {      // а если количество слушающих клиентов = 0 то..
						s.DeleteRoom(roomID) // удаляем комнату
					}
				}
			}

		}
	}
	return nil
}

func (s *Subix) Unsub(sessionKey Key) {
	s.deleteClient(sessionKey)
}
