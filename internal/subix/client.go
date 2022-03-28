package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
	"time"
)

type Employees map[ID]*Employee

type Employee struct {
	EmpID   int
	rooms   Rooms
	clients Clients
}

type Clients map[Key]*Client
type ClientsWithEvents map[Key]*ClientWithEvents

type Client struct {
	EmployeeID int
	//ExpectedEvents   map[model.EventType]bool
	Ch               chan *model.SubscriptionBody
	task             *scheduler.Task
	sessionExpiresAt int64
	sessionKey       Key
	marked           bool
}

type EventCollection map[model.EventType]bool
type ClientWithEvents struct {
	Client *Client
	Events EventCollection
}

func (s *Subix) CreateEmployeeIfNotExists(empID int) *Employee {
	emp, ok := s.employees[empID]
	if !ok {
		emp = &Employee{
			EmpID:   empID,
			rooms:   Rooms{},
			clients: Clients{},
		}
		s.employees[empID] = emp
	}
	return emp
}

func (s *Subix) deleteEmployee(empID int) {
	emp, ok := s.employees[empID]
	if ok { // если вдруг не удается найти, то просто пропускаем
		delete(s.employees, empID)           // удаление из глобальной мапы
		for _, client := range emp.clients { // определяем тех клиентов которых надо удалить из глобальной мапы
			delete(s.clients, client.sessionKey) // удаление
		}

		for _, room := range emp.rooms { // а здесь определяем мемберов, к которым относятся к пользователю

			delete(room.Empls, empID)

			if len(room.Empls) == 0 {
				s.DeleteRoom(room.RoomID) // удаляем по отдельности через функцию
			} else {
				for _, client := range emp.clients {
					delete(room.clientsWithEvents, client.sessionKey)
					if len(room.clientsWithEvents) == 0 {
						s.DeleteRoom(room.RoomID)
					}
				}
			}

		}
		emp.clients = nil
		//emp.rooms = nil // на всякий случай заnullяем мапу
		// теперь на этого пользователя не должно остаться ссылок как и на его клиентов
	}

}

func (s *Subix) deleteClient(sessionKey Key) error {
	client, ok := s.clients[sessionKey]
	if ok {
		delete(s.clients, client.sessionKey)
		err := s.sched.DropTask(&client.task)
		if err != nil {
			return cerrors.Wrap(err, "не удалось удалить клиента")
		}
		select {
		case x, ok := <-client.Ch:
			if ok {
				select {
				case client.Ch <- x:
				default:
				}
				close(client.Ch)
			}
		default:
			close(client.Ch)
		}

		emp, ok := s.employees[client.EmployeeID]
		if ok {
			delete(emp.clients, client.sessionKey)
			if len(emp.clients) == 0 {
				s.deleteEmployee(emp.EmpID)
			}
		}

		for _, room := range emp.rooms {
			delete(room.clientsWithEvents, client.sessionKey)
			if len(room.clientsWithEvents) == 0 {
				s.DeleteRoom(room.RoomID)
			}
		}
	}
	//println("удален клиент", client.sessionKey)
	return nil
}

func (s *Subix) scheduleMarkClient(client *Client, expAt int64) (err error) {
	client.task, err = s.sched.AddTask(
		func() {
			eventBody := &model.TokenExpired{
				Message: "используйте mutation.RefreshTokens для того чтобы возобновить получение данных, иначе соединение закроется",
			}
			s.writeToClient(
				client,
				&model.SubscriptionBody{
					Event: getEventTypeByEventResult(eventBody),
					Body:  eventBody,
				},
			)
			client.marked = true // теперь будем знать что этому клиенту не надо отправлять события
			//println("токен клиента истек, помечаю клиента", client)
			err := s.scheduleExpiredClient(client)
			if err != nil {
				panic(err)
			}

		},
		expAt,
	)

	return err
}

func (s *Subix) scheduleExpiredClient(client *Client) (err error) {
	client.task, err = s.sched.AddTask(
		func() {
			s.deleteClient(client.sessionKey) // клиент не обновил токен, удаляем его
		},
		time.Now().Unix()+rules.LifetimeOfMarkedClient,
	)

	return err
}

func (s *Subix) ExtendClientSession(sessionKey Key, expAt int64) (err error) {
	client, ok := s.clients[sessionKey]
	if !ok {
		return cerrors.New("не удалось продлить сессию, клиент не найден")
	}
	err = s.sched.DropTask(&client.task)
	if err != nil {
		return err
	}
	err = s.scheduleMarkClient(client, expAt)
	if err != nil {
		return err
	}
	client.marked = false
	// сессия успешно продлена
	return nil
}

func (s Subix) ClientCollection(sessionKey Key) (collection []*model.ListenedChat) {
	client, _ := s.clients[sessionKey]                          // предполагается что сессия с таким ключом существует
	for _, room := range s.employees[client.EmployeeID].rooms { // по комнатам пользователя
		listenedChat := &model.ListenedChat{ID: room.RoomID}
		for event := range room.clientsWithEvents[sessionKey].Events {
			listenedChat.Events = append(listenedChat.Events, event)
		}
		collection = append(collection, listenedChat)
	}
	return collection
}
