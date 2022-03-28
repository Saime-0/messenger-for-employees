package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type (
	ID  = int
	Key = string // (sessionKey | sessionKey) len = 20, any symbols "Twenty-Digit-Session-Key": "[.]20"
)

func (s *Subix) NotifyChatMembers(chat ID, body model.EventResult) {
	err := s.spam(
		[]ID{chat},
		s.repo.Subscribers.Members,
		body,
	)
	if err != nil {
		panic(err)
	}
}
func (s *Subix) NotifyChats(chats []ID, body model.EventResult) {
	err := s.spam(
		chats,
		s.repo.Subscribers.Members,
		body,
	)
	if err != nil {
		panic(err)
	}
}
func (s *Subix) NotifyRoomReaders(room ID, body model.EventResult) error {
	err := s.spam(
		[]ID{room},
		s.repo.Subscribers.RoomReaders,
		body,
	)
	if err != nil {
		return cerrors.Wrap(err, "ивент небыл разослан")
	}
	return nil
}

func (s *Subix) spam(objects []ID, meth repository.QueryUserGroup, body interface{}) error {

	switch body.(type) {

	case *model.DeleteMember: // чтобы участник(пользователь) перестал получать события
		s.DeleteRoom(body.(*model.DeleteMember).ID)

	case *model.NewMessage: // ожидается что в objects будут RoomID комнат
		readers, err := s.repo.Subscribers.RoomReaders(objects)
		if err != nil {
			return cerrors.Wrap(err, "не найдены участники комнаты")
		}
		s.informMembers(readers, body.(model.EventResult))
		return nil
	}

	s.informChat(objects, body.(model.EventResult))
	return nil
}
