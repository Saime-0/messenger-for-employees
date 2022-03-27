package piper

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/internal/validator"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

//func (n Node) ChatExists(chatID int) (fail bool) {
//	n.SwitchMethod("ChatExists", &bson.M{
//		"chatID": chatID,
//	})
//	defer n.MethodTiming()
//
//	exists, err := n.Dataloader.UnitExistsByID(chatID, model.UnitTypeChat)
//	if err != nil {
//		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
//		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
//		return true
//	}
//	if !exists {
//		n.SetError(resp.ErrBadRequest, fmt.Sprintf("chat(id:%d) is not exists", chatID))
//		return true
//	}
//	return
//}

//func (n Node) UserExists(userID int) (fail bool) {
//	n.SwitchMethod("UserExists", &bson.M{
//		"userID": userID,
//	})
//	defer n.MethodTiming()
//
//	exists, err := n.Dataloader.UnitExistsByID(userID, model.UnitTypeUser)
//	if err != nil {
//		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
//		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
//		return true
//	}
//	if !exists {
//		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is not exists", userID))
//		return true
//	}
//	return
//}

// ValidParams
//  with side effect
func (n Node) ValidParams(params **model.Params) (fail bool) {
	n.SwitchMethod("ValidParams", &bson.M{
		"params": params,
	})
	defer n.MethodTiming()

	if *params == nil {
		*params = &model.Params{
			Limit:  kit.IntPtr(rules.MaxLimit), // ! unsafe
			Offset: kit.IntPtr(0),
		}
		return
	}
	if (*params).Limit != nil {
		if !validator.ValidateLimit(*(*params).Limit) {
			n.SetError(resp.ErrBadRequest, "невалидное значение лимита")
			return true
		}
	} else {
		(*params).Limit = kit.IntPtr(rules.MaxLimit)
	}
	if (*params).Offset != nil {
		if !validator.ValidateOffset(*(*params).Offset) {
			n.SetError(resp.ErrBadRequest, "невалидное значение смещения")
			return true
		}
	} else {
		(*params).Offset = kit.IntPtr(0)
	}

	return
}

func (n Node) ValidNameFragment(fragment string) (fail bool) {
	n.SwitchMethod("ValidNameFragment", &bson.M{
		"fragment": fragment,
	})
	defer n.MethodTiming()

	if !(validator.ValidateEmployeeFullName(fragment) || validator.ValidatePartOfName(fragment)) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для фрагмента имени")
		return true
	}
	return
}

func (n *Node) ValidNote(note string) (fail bool) {
	n.SwitchMethod("ValidNote", &bson.M{
		"note": note,
	})
	defer n.MethodTiming()

	if !validator.ValidateNote(note) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для заметки")
		return true
	}
	return
}

func (n Node) ValidID(id int) (fail bool) {
	n.SwitchMethod("ValidID", &bson.M{
		"id": id,
	})
	defer n.MethodTiming()

	if !validator.ValidateID(id) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

func (n Node) ValidParentRoomID(id, parent int) (fail bool) {
	n.SwitchMethod("ValidParentRoomID", &bson.M{
		"id":     id,
		"parent": parent,
	})
	defer n.MethodTiming()

	if !validator.ValidateID(parent) || id == parent {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

//func (n Node) ValidFindMessagesInRoom(find *model.FindMessagesInRoom) (fail bool) {
//	n.SwitchMethod("ValidFindMessagesInRoom", &bson.M{
//		"find": find,
//	})
//	defer n.MethodTiming()
//
//	if find.Count <= 0 ||
//		find.Count > *n.cfg.MaximumNumberOfMessagesPerRequest ||
//		find.Created == model.MessagesCreatedBefore && find.StartMessageID-find.Count < 0 {
//		n.SetError(resp.ErrBadRequest, "неверное значение количества сообщений")
//		return true
//	}
//	if !validator.ValidateID(find.StartMessageID) {
//		n.SetError(resp.ErrBadRequest, "неверный ID сообщения")
//		return true
//	}
//
//	return
//}

func (n Node) ValidPassword(password string) (fail bool) {
	n.SwitchMethod("ValidPassword", &bson.M{
		"password": password,
	})
	defer n.MethodTiming()

	if !validator.ValidatePassword(password) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для пароля")
		return true
	}
	return
}

func (n Node) ValidDomain(domain string) (fail bool) {
	n.SwitchMethod("ValidDomain", &bson.M{
		"domain": domain,
	})
	defer n.MethodTiming()

	if !validator.ValidateDomain(domain) {
		n.SetError(resp.ErrBadRequest, "невалидный домен")
		return true
	}
	return
}

func (n Node) ValidRoomName(name string) (fail bool) {
	n.SwitchMethod("ValidRoomName", &bson.M{
		"name": name,
	})
	defer n.MethodTiming()

	if !validator.ValidateRoomName(name) {
		n.SetError(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (n Node) ValidEmployeePartOfName(partOfName string) (fail bool) {
	n.SwitchMethod("ValidEmployeePartOfName", &bson.M{
		"partOfName": partOfName,
	})
	defer n.MethodTiming()

	if !validator.ValidatePartOfName(partOfName) {
		n.SetError(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (n Node) ValidEmployeeFullName(fullName string) (fail bool) {
	n.SwitchMethod("ValidEmployeeFullName", &bson.M{
		"fullName": fullName,
	})
	defer n.MethodTiming()

	if !validator.ValidateEmployeeFullName(fullName) {
		n.SetError(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (n Node) DomainIsFree(domain string) (fail bool) {
	n.SwitchMethod("DomainIsFree", &bson.M{
		"domain": domain,
	})
	defer n.MethodTiming()

	free, err := n.repos.Units.DomainIsFree(domain)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if !free {
		n.SetError(resp.ErrBadRequest, "домен занят")
		return true
	}
	return
}

func (n Node) EmailIsFree(email string) (fail bool) {
	n.SwitchMethod("EmailIsFree", &bson.M{
		"email": email,
	})
	defer n.MethodTiming()

	free, err := n.repos.Users.EmailIsFree(email)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if !free {
		n.SetError(resp.ErrBadRequest, "такая почта уже занята кем-то")
		return true
	}
	return
}

// IsMember does not need if the Can method is used..
func (n Node) IsMember(userID, roomID int) (fail bool) {
	n.SwitchMethod("IsMember", &bson.M{
		"userID": userID,
		"roomID": roomID,
	})
	defer n.MethodTiming()

	isMember, err := n.Dataloader.EmployeeIsRoomMember(userID, roomID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if !isMember {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is not member of room(id:%d)", userID, roomID))
		return true
	}

	return
}

func (n Node) IsNotMember(userID, chatID int) (fail bool) {
	n.SwitchMethod("IsNotMember", &bson.M{
		"userID": userID,
		"chatID": chatID,
	})
	defer n.MethodTiming()

	isMember, err := n.Dataloader.EmployeeIsRoomMember(userID, chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if isMember {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is member of chat(id:%d)", userID, chatID))
		return true
	}
	return
}

func (n Node) MessageExists(roomID, msgID int) (fail bool) {
	n.SwitchMethod("MessageExists", &bson.M{
		"roomID": roomID,
		"msgID":  msgID,
	})
	defer n.MethodTiming()

	exists, err := n.Dataloader.MessageExists(roomID, msgID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if exists {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("message(id:%d) is not exists in room(id:%d)", roomID, msgID))
		return true
	}
	return
}

func (n Node) RolesLimit(chatID int) (fail bool) {
	n.SwitchMethod("RolesLimit", &bson.M{
		"chatID": chatID,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.GetCountChatRoles(chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= *n.cfg.MaxRolesInChat {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("limit on the number of chat roles has been reached (MaxRolesInChat = %d)", *n.cfg.MaxRolesInChat))
		return true
	}
	return
}

func (n Node) RoomsLimit(chatID int) (fail bool) {
	n.SwitchMethod("RoomsLimit", &bson.M{
		"chatID": chatID,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.GetCountRooms(chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= *n.cfg.MaxCountRooms {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("limit of the number of rooms in the chat has been reached (MaxCountRooms = %d)", *n.cfg.MaxCountRooms))
		return true
	}
	return
}

func (n Node) RoomExists(roomID int) (fail bool) {
	n.SwitchMethod("RoomExists", &bson.M{
		"roomID": roomID,
	})
	defer n.MethodTiming()

	//if !n.repos.Rooms.RoomExistsByID(roomID) {
	exists, err := n.Dataloader.RoomExistsByID(roomID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if !exists {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("room(id:%d) is not exists", roomID))
		return true
	}
	return
}

func (n Node) RoleExists(chatID, roleID int) (fail bool) {
	n.SwitchMethod("RoleExists", &bson.M{
		"chatID": chatID,
		"roleID": roleID,
	})
	defer n.MethodTiming()

	exists, err := n.repos.Chats.RoleExistsByID(chatID, roleID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if !exists {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("role(id:%d) is not exists", roleID))
		return true
	}
	return
}

//func (n Node) MembersLimit(chatID int) (fail bool) {
//	n.SwitchMethod("MembersLimit", &bson.M{
//		"chatID": chatID,
//	})
//	defer n.MethodTiming()
//
//	count, err := n.repos.Chats.CountMembers(chatID)
//	if err != nil {
//		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
//		return true
//	}
//	if count >= *n.cfg.MaxMembersOnChat {
//		n.SetError(resp.ErrBadRequest, fmt.Sprintf("limit of the number of participants in the chat has been reached (MaxMembersOnChat = %d)", *n.cfg.MaxMembersOnChat))
//		return true
//	}
//	return
//}

func (n Node) ChatIsNotPrivate(chatID int) (fail bool) {
	n.SwitchMethod("ChatIsNotPrivate", &bson.M{
		"chatID": chatID,
	})
	defer n.MethodTiming()

	if n.repos.Chats.ChatIsPrivate(chatID) {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("chat(id:%d) is private", chatID))
		return true
	}
	return
}

func (n Node) UserExistsByRequisites(input *models.LoginRequisites) (fail bool) {
	n.SwitchMethod("UserExistsByRequisites", &bson.M{
		"input": input,
	})
	defer n.MethodTiming()

	exists, err := n.repos.Users.UserExistsByRequisites(input)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if !exists {
		n.SetError(resp.ErrBadRequest, "неверный логин или пароль ")
		return true
	}
	return
}

func (n Node) GetUserIDByRequisites(input *models.LoginRequisites, userID *int) (fail bool) {
	n.SwitchMethod("GetUserIDByRequisites", &bson.M{
		"input":  input,
		"userID": userID,
	})
	defer n.MethodTiming()

	_uid, err := n.repos.Users.GetUserIdByRequisites(input)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*userID = _uid
	return
}

func (n Node) IsNotBanned(userID, chatID int) (fail bool) {
	n.SwitchMethod("IsNotBanned", &bson.M{
		"userID": userID,
		"chatID": chatID,
	})
	defer n.MethodTiming()

	banned, err := n.repos.Chats.UserIsBanned(userID, chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}

	if banned {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is banned in chat(id:%d)", userID, chatID))
		return true
	}
	return
}

func (n Node) IsBanned(userID, chatID int) (fail bool) {
	n.SwitchMethod("IsBanned", &bson.M{
		"userID": userID,
		"chatID": chatID,
	})
	defer n.MethodTiming()

	banned, err := n.repos.Chats.UserIsBanned(userID, chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}

	if !banned {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is not banned in chat(id:%d)", userID, chatID))
		return true
	}
	return
}

func (n Node) GetChatIDByRoom(roomID int, chatID *int) (fail bool) {
	n.SwitchMethod("GetChatIDByRoom", &bson.M{
		"roomID": roomID,
		"chatID": chatID,
	})
	defer n.MethodTiming()

	_chatID, err := n.repos.Rooms.GetChatIDByRoomID(roomID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, fmt.Sprintf("room(id:%d) does not apply to any chat", roomID))
		return true
	}
	*chatID = _chatID
	return
}

func (n Node) GetChatIDByAllow(allowID int, chatID *int) (fail bool) {
	n.SwitchMethod("GetChatIDByAllow", &bson.M{
		"allowID": allowID,
		"chatID":  chatID,
	})
	defer n.MethodTiming()

	_chatID, err := n.repos.Rooms.GetChatIDByAllowID(allowID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if _chatID == 0 {
		n.SetError(resp.ErrInternalServerError, fmt.Sprintf("allow(id:%d) does not apply to any chat", allowID))
		return true
	}
	*chatID = _chatID
	return
}

func (n Node) GetChatIDByMember(memberID int, chatID *int) (fail bool) {
	n.SwitchMethod("GetChatIDByMember", &bson.M{
		"memberID": memberID,
		"chatID":   chatID,
	})
	defer n.MethodTiming()

	var err error
	*chatID, err = n.Dataloader.ChatIDByMemberID(memberID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if *chatID < 1 {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("member(id:%d) does not apply to any chat", memberID))
		return true
	}

	return
}

func (n Node) GetMemberBy(userID, chatID int, memberID *int) (fail bool) {
	n.SwitchMethod("GetMemberBy", &bson.M{
		"userID":   userID,
		"chatID":   chatID,
		"memberID": memberID,
	})
	defer n.MethodTiming()

	by, err := n.Dataloader.FindMemberBy(userID, chatID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if by == nil || *by == 0 {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is not member of chat(id:%d)", userID, chatID))
		return true
	}
	*memberID = *by
	return
}

func (n Node) IsNotMuted(memberID int) (fail bool) {
	n.SwitchMethod("IsNotMuted", &bson.M{
		"memberID": memberID,
	})
	defer n.MethodTiming()

	muted, err := n.repos.Chats.MemberIsMuted(memberID)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if muted {
		n.SetError(resp.ErrBadRequest, "участник заглушен")
		return true
	}
	return
}

func (n Node) ValidSessionKey(sessionKey string) (fail bool) {
	n.SwitchMethod("ValidSessionKey", &bson.M{
		"sessionKey": sessionKey,
	})
	defer n.MethodTiming()

	if !validator.ValidateSessionKey(sessionKey) {
		n.SetError(resp.ErrBadRequest, "невалидный ключ сессии. Требования: Length 20 characters & (Upper or lower case letters | Special symbols (-.=)")
		return true
	}
	return
}

func (n Node) ValidEmail(email string) (fail bool) {
	n.SwitchMethod("ValidateEmail", &bson.M{
		"email": email,
	})
	defer n.MethodTiming()

	if !validator.ValidateEmail(email) {
		n.SetError(resp.ErrBadRequest, "невалидный email")
		return true
	}
	return
}

func (n Node) GetRegistrationSession(email, code string, regi **models.RegisterData) (fail bool) {
	n.SwitchMethod("GetRegistrationSession", &bson.M{
		"email": email,
		"code":  code,
		"regi":  regi,
	})
	defer n.MethodTiming()

	var err error
	*regi, err = n.repos.Users.GetRegistrationSession(email, code)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if *regi == nil {
		n.SetError(resp.ErrBadRequest, "сессии не существует")
		return true
	}
	return
}

func (n Node) ValidRegisterInput(input *model.RegisterInput) (fail bool) {
	n.SwitchMethod("ValidRegisterInput", &bson.M{
		"input": input,
	})
	defer n.MethodTiming()

	switch {
	case !validator.ValidateDomain(input.Domain):
		n.SetError(resp.ErrBadRequest, "домен не соответствует требованиям")
		return true

	case !validator.ValidateEmployeeFullName(input.Name):
		n.SetError(resp.ErrBadRequest, "имя не соответствует требованиям")
		return true

	case !validator.ValidateEmail(input.Email):
		n.SetError(resp.ErrBadRequest, "имеил не соответствует требованиям")
		return true

	case !validator.ValidatePassword(input.Password):
		n.SetError(resp.ErrBadRequest, "пароль не соответствует требованиям")
		return true

	}

	return
}

func (n Node) UserHasAccessToChats(userID int, chats *[]int, submembers **[]*models.SubUser) (fail bool) {
	n.SwitchMethod("UserHasAccessToChats", &bson.M{
		"userID":     userID,
		"chats":      chats,
		"submembers": submembers,
	})
	defer n.MethodTiming()

	if !validator.ValidateIDs(*chats) {
		n.SetError(resp.ErrBadRequest, "chatID is not valid")
		return true
	}
	members, noAccessTo, err := n.repos.Chats.UserHasAccessToChats(userID, chats)
	if err != nil {
		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if noAccessTo != 0 {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) does not have access to chat(id:%d)", userID, noAccessTo))
		return true
	}
	*submembers = &members
	return
}
