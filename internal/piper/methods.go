package piper

import (
	"fmt"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/models"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/rules"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"github.com/saime-0/messenger-for-employee/internal/validator"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

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

func (n Node) ValidMsgCount(count int) (fail bool) {
	n.SwitchMethod("ValidMsgCount", &bson.M{
		"count": count,
	})
	defer n.MethodTiming()

	if !(count > 0 && count <= rules.MaxMsgCount) {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("count must be more then 0 or less or equal then %d", rules.MaxMsgCount))
		return true
	}
	return
}

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

func (n Node) EmailIsFree(email string) (fail bool) {
	n.SwitchMethod("EmailIsFree", &bson.M{
		"email": email,
	})
	defer n.MethodTiming()

	free, err := n.repos.Employees.EmailIsFree(email)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
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
func (n Node) IsMember(employeeID, roomID int) (fail bool) {
	n.SwitchMethod("IsMember", &bson.M{
		"employeeID": employeeID,
		"roomID":     roomID,
	})
	defer n.MethodTiming()

	isMember, err := n.Dataloader.EmployeeIsRoomMember(employeeID, roomID)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if !isMember {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("user(id:%d) is not member of room(id:%d)", employeeID, roomID))
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

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if exists {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("message(id:%d) is not exists in room(id:%d)", roomID, msgID))
		return true
	}
	return
}

func (n Node) RoomExists(roomID int) (fail bool) {
	n.SwitchMethod("RoomExists", &bson.M{
		"roomID": roomID,
	})
	defer n.MethodTiming()

	//if !n.repos.EmployeeRooms.RoomExistsByID(roomID) {
	exists, err := n.Dataloader.RoomExistsByID(roomID)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrBadRequest, "произошла ошибка во время обработки данных")
		return true
	}
	if !exists {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("room(id:%d) is not exists", roomID))
		return true
	}
	return
}

func (n Node) UserExistsByRequisites(input *models.LoginRequisites) (fail bool) {
	n.SwitchMethod("EmployeeExistsByRequisites", &bson.M{
		"input": input,
	})
	defer n.MethodTiming()

	exists, err := n.repos.Employees.EmployeeExistsByRequisites(input)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if !exists {
		n.SetError(resp.ErrBadRequest, "неверный логин или пароль ")
		return true
	}
	return
}

func (n Node) GetEmployeeIDByRequisites(input *models.LoginRequisites, employeeID *int) (fail bool) {
	n.SwitchMethod("GetEmployeeIDByRequisites", &bson.M{
		"input":      input,
		"employeeID": employeeID,
	})
	defer n.MethodTiming()

	_uid, err := n.repos.Employees.GetEmployeeIDByRequisites(input)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*employeeID = _uid
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

func (n Node) ValidMessageBody(body *string) (fail bool) {
	n.SwitchMethod("ValidMessageBody", &bson.M{
		"body": body,
	})
	defer n.MethodTiming()

	if ok, reason := validator.ValidateMessageBody(body); !ok {
		n.SetError(resp.ErrBadRequest, reason)
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

func (n Node) EmployeeHasAccessToRooms(employeeID int, roomIDs []int) (fail bool) {
	n.SwitchMethod("EmployeeHasAccessToRooms", &bson.M{
		"employeeID": employeeID,
		"roomIDs":    roomIDs,
	})
	defer n.MethodTiming()

	if !validator.ValidateIDs(roomIDs) {
		n.SetError(resp.ErrBadRequest, "roomID is not valid")
		return true
	}
	noAccessTo, err := n.repos.Rooms.EmployeeHasAccessToRooms(employeeID, roomIDs)
	if err != nil {

		n.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if noAccessTo != 0 {
		n.SetError(resp.ErrBadRequest, fmt.Sprintf("employee(id:%d) does not have access to room(id:%d)", employeeID, noAccessTo))
		return true
	}
	return false
}
