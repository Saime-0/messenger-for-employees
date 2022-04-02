package subix

import (
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"log"
)

func (s *Subix) writeToRoom(ignoreEventCollection bool, body model.EventResult, roomIDs ...int) {
	for _, roomID := range roomIDs {
		room, ok := s.rooms[roomID]
		if !ok {
			continue
		}
		s.writeToClientsWithEvents(
			ignoreEventCollection,
			room.clientsWithEvents,
			body,
			getEventTypeByEventResult(body),
		)

	}
}

func (s *Subix) writeToEmployees(body model.EventResult, empIDs ...int) {
	eventType := getEventTypeByEventResult(body)
	for _, empID := range empIDs {
		emp, ok := s.employees[empID]
		if !ok {
			continue
		}

		for _, client := range emp.clients {
			// У клиента нет набора ивентов которые он хочет получать
			// они есть только у участников, тк для каждой отдельной комнаты можно установить свою коллекцию,
			// а остальные ивенты, не связанные с комнатами, клиенты получают обязательно (прим. R)
			s.writeToClient(
				client,
				&model.SubscriptionBody{
					Event: eventType,
					Body:  body,
				},
			)
		}

	}
}

func (s *Subix) writeToAllEmployees(body model.EventResult) {
	eventType := getEventTypeByEventResult(body)
	for _, emp := range s.employees {

		for _, client := range emp.clients {
			s.writeToClient(
				client,
				&model.SubscriptionBody{
					Event: eventType,
					Body:  body,
				},
			)
		}

	}
}

func (s *Subix) writeToClientsWithEvents(ignoreEventCollection bool, clientsWithEvents ClientsWithEvents, body model.EventResult, eventType model.EventType) {
	for _, clientWithEvents := range clientsWithEvents {
		if !ignoreEventCollection {
			if _, ok := clientWithEvents.Events[eventType]; !ok { // если он не слушает эти события, то...
				continue // ...и слать их ему не надо, просто пропускаем этого клиента
			}
		}

		s.writeToClient(
			clientWithEvents.Client,
			&model.SubscriptionBody{
				Event: eventType,
				Body:  body,
			},
		)
	}
}

func (s *Subix) writeToClient(client *Client, subbody *model.SubscriptionBody) {
	log.Printf("%s client: %#v", utils.GetCallerPos(), client) // debug
	if (*client).marked {
		return
	}
	select {
	case (*client).Ch <- subbody: // success
	default: // канал никто не слушает
		if client != nil {
			defer s.deleteClient(client.sessionKey)
		}
	}
}

func getEventTypeByEventResult(body model.EventResult) model.EventType {
	switch body.(type) {
	case *model.NewMessage:
		return model.EventTypeNewMessage
	case *model.DropTag:
		return model.EventTypeDropTag
	case *model.EmpTagAction:
		return model.EventTypeEmpTagAction
	case *model.MemberAction:
		return model.EventTypeMemberAction
	case *model.DropRoom:
		return model.EventTypeDropRoom
	case *model.TokenExpired:
		return model.EventTypeTokenExpired
	default:
		panic("not implemented")
	}
}
