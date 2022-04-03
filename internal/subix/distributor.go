package subix

import (
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
)

type (
	ID  = int
	Key = string // (sessionKey | sessionKey) len = 20, any symbols "Twenty-Digit-Session-Key": "[.]20"
)

// NotifyRoomMembers клиенты получат только те ивенты на которые они подписались
func (s *Subix) NotifyRoomMembers(body model.EventResult, rooms ...ID) error {

	s.writeToRoom(false, body.(model.EventResult), rooms...)

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

// все учсастники(клиенты) получат ивент
func (s *Subix) NotifyAllRoomMembers(body model.EventResult, rooms ...ID) error {

	s.writeToRoom(true, body.(model.EventResult), rooms...)

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

// можно отправлять любой ивент, он все равно будет отправлен клиенту
func (s *Subix) NotifyEmployees(body model.EventResult, empIDs ...ID) error {

	s.writeToEmployees(body.(model.EventResult), empIDs...)

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

// отправить всем клиентам
func (s *Subix) NotifyAllEmployees(body model.EventResult) error {

	s.writeToAllEmployees(body.(model.EventResult))

	if err := s.handleEventTypeAfterNotify(body); err != nil {
		return cerrors.Wrap(err, "handleEventTypeAfterNotify failure")
	}
	return nil
}

func (s *Subix) handleEventTypeAfterNotify(body interface{}) error {

	switch body.(type) {

	case *model.DropRoom: // чтобы участники перестали получать события
		s.DeleteRoom(body.(*model.DropRoom).RoomID)

	case *model.MemberAction:
		if body.(*model.MemberAction).Action == model.ActionDel {
			for _, roomID := range body.(*model.MemberAction).RoomIDs {
				room, ok := s.rooms[roomID]
				if ok {
					delete(room.Empls, body.(*model.MemberAction).EmpID)
				}
			}
		}

	}

	return nil
}
