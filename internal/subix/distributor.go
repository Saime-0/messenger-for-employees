package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
)

type (
	ID  = int
	Key = string // (sessionKey | sessionKey) len = 20, any symbols "Twenty-Digit-Session-Key": "[.]20"
)

// NotifyRoomMembers клиенты получат только те ивенты на которые они подписались
func (s *Subix) NotifyRoomMembers(body model.EventResult, rooms ...ID) error {

	s.writeToRoom(body.(model.EventResult), rooms...)

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

// NotifyEmployees можно отправлять любой ивент, он все равно будет отправлен клиенту
func (s *Subix) NotifyEmployees(body model.EventResult, empIDs ...ID) error {

	s.writeToEmployees(body.(model.EventResult), empIDs...)

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

func (s *Subix) handleEventTypeAfterNotify(body interface{}) error {

	switch body.(type) {

	case *model.DeleteRoom: // чтобы участники перестали получать события
		for _, id := range body.(*model.DeleteRoom).RoomsID {
			s.DeleteRoom(id)
		}

	case *model.RemoveMember:
		for _, roomID := range body.(*model.RemoveMember).RoomIDs {
			delete(s.rooms[roomID].Empls, body.(*model.RemoveMember).EmpID)
		}
	}

	return nil
}
