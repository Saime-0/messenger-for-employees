package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"strings"
)

func (s *Subix) writeToRoom(body model.EventResult, roomIDs ...int) {
	for _, roomID := range roomIDs {
		room, ok := s.rooms[roomID]
		if !ok {
			continue
		}
		s.writeToClientsWithEvents(
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
			// у клиента нет набора ивентов которые он хочет получать
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

func (s *Subix) writeToClientsWithEvents(clientsWithEvents ClientsWithEvents, body model.EventResult, eventType model.EventType) {
	for _, clientWithEvents := range clientsWithEvents {
		if _, ok := clientWithEvents.Events[eventType]; !ok { // если он не слушает эти события, то..
			continue // ..и слать их ему не надо, просто скипаем этого клиента
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

func getEventType(body model.EventResult) string {
	bodyType := fmt.Sprintf("%T", body)
	dot := strings.LastIndex(
		bodyType,
		".",
	)
	if dot == -1 {
		panic("invalid index")
	}
	return strings.ToUpper(bodyType[dot+1:])
}

func getEventTypeByEventResult(body model.EventResult) model.EventType {
	switch body.(type) {
	case *model.NewMessage:
		return model.EventTypeNewMessage
	case *model.UpdateEmpFirstName:
		return model.EventTypeUpdateEmpFirstName
	case *model.UpdateEmpLastName:
		return model.EventTypeUpdateEmpLastName
	case *model.GiveTagToEmp:
		return model.EventTypeGiveTagToEmp
	case *model.TakeTagFromEmp:
		return model.EventTypeTakeTagFromEmp
	case *model.RemoveTagFromEmp:
		return model.EventTypeRemoveTagFromEmp
	case *model.NewMember:
		return model.EventTypeNewMember
	case *model.RemoveMember:
		return model.EventTypeRemoveMember
	case *model.CreateTag:
		return model.EventTypeCreateTag
	case *model.UpdateTag:
		return model.EventTypeUpdateTag
	case *model.DeleteTag:
		return model.EventTypeDeleteTag
	case *model.UpdateRoomName:
		return model.EventTypeUpdateRoomName
	case *model.DeleteRoom:
		return model.EventTypeDeleteRoom
	case *model.TokenExpired:
		return model.EventTypeTokenExpired
	default:
		panic("not implemented")
	}
}
