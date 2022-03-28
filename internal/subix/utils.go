package subix

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"strings"
)

func (s *Subix) informMembers(membersID []int, body model.EventResult) {
	for _, memberID := range membersID {
		member, ok := s.rooms[memberID]
		if !ok {
			continue
		}
		s.writeToMembers(
			member.clientsWithEvents,
			body,
			getEventTypeByEventResult(body),
		)

	}
}

func (s *Subix) informChat(chatsID []int, body model.EventResult) {
	for _, chatID := range chatsID {
		chat, ok := s.chats[chatID]
		if !ok {
			continue
		}

		for _, member := range chat.members {

			s.writeToMembers(
				member.clientsWithEvents,
				body,
				getEventTypeByEventResult(body),
			)

		}

	}
}

// deprecated
func (s *Subix) writeToUsers(usersID []int, body model.EventResult) {
	eventType := getEventTypeByEventResult(body)
	for _, userID := range usersID {
		user, ok := s.employees[userID]
		if !ok {
			continue
		}
		for _, client := range user.clients {
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

func (s *Subix) writeToMembers(clientsWithEvents ClientsWithEvents, body model.EventResult, eventType model.EventType) {
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
	default: // client chan is close
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
