package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/messenger-for-employee/graph/generated"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *employeeResolver) Tags(ctx context.Context, obj *model.Employee) (*model.Tags, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("EmpID.Tags", &bson.M{
		"employeeID (obj.EmpID)": obj.EmpID,
	})
	defer node.MethodTiming()

	rooms, err := r.Dataloader.Tags(obj.EmpID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("ошибка при попытке получить данные")
	}

	return rooms, nil
}

func (r *listenCollectionResolver) Collection(ctx context.Context, obj *model.ListenCollection) ([]*model.ListenedChat, error) {
	collection := r.Subix.ClientCollection(obj.SessionKey)
	return collection, nil
}

func (r *meResolver) Rooms(ctx context.Context, obj *model.Me, params model.Params) (*model.Rooms, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Me.EmployeeRooms", &bson.M{
		"employeeID (obj.EmpID.EmpID)": obj.Employee.EmpID,
		"params":                       params,
	})
	defer node.MethodTiming()

	rooms, err := r.Dataloader.Rooms(obj.Employee.EmpID, &params)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("ошибка при попытке получить данные")
	}

	return rooms, nil
}

func (r *memberResolver) Employee(ctx context.Context, obj *model.Member) (*model.Employee, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.EmpID", &bson.M{
		"employeeID (obj.EmpID.EmpID)": obj.Employee.EmpID,
	})
	defer node.MethodTiming()

	employee, err := r.Dataloader.Employee(obj.Employee.EmpID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных")
	}

	return employee, nil
}

func (r *memberResolver) Room(ctx context.Context, obj *model.Member) (*model.Room, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Member.Room", &bson.M{
		"roomID (obj.Room.Rooms)": obj.Room.RoomID,
	})
	defer node.MethodTiming()

	var clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID

	room, err := r.Dataloader.Room(clientID, obj.Room.RoomID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("ошибка при попытке получить данные")
	}

	return room, nil
}

func (r *messageResolver) Room(ctx context.Context, obj *model.Message) (*model.Room, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.Room", &bson.M{
		"roomID (obj.Room.Rooms)": obj.Room.RoomID,
	})
	defer node.MethodTiming()

	var clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID

	room, err := r.Dataloader.Room(clientID, obj.Room.RoomID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("ошибка при попытке получить данные")
	}

	return room, nil
}

func (r *messageResolver) Employee(ctx context.Context, obj *model.Message) (*model.Employee, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.EmpID", &bson.M{
		"employeeID (obj.EmpID.EmpID)": obj.Employee.EmpID,
	})
	defer node.MethodTiming()

	employee, err := r.Dataloader.Employee(obj.Employee.EmpID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных")
	}

	return employee, nil
}

func (r *messageResolver) TargetMsgID(ctx context.Context, obj *model.Message) (*model.Message, error) {
	if obj.TargetMsgID == nil {
		return nil, nil
	}
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.TargetMsgID", &bson.M{
		"msgID (obj.TargetMsgID.MsgID)":       obj.TargetMsgID.MsgID,
		"msgID (obj.TargetMsgID.Room.RoomID)": obj.TargetMsgID.Room.RoomID,
	})
	defer node.MethodTiming()

	message, err := r.Dataloader.Message(obj.TargetMsgID.Room.RoomID, obj.TargetMsgID.MsgID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных")
	}
	return message, nil
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room, startMsg int, created model.MsgCreated, count int) (*model.Messages, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Messages", &bson.M{
		"roomID (obj.Rooms)": obj.RoomID,
		"startMsg":           startMsg,
		"created":            created,
		"count":              count,
	})
	defer node.MethodTiming()

	if node.ValidID(startMsg) ||
		node.ValidMsgCount(count) {
		return nil, cerrors.New(node.GetError().Error)
	}

	messages, err := r.Services.Repos.Rooms.RoomMessages(obj.RoomID, startMsg, created, count)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных")
	}

	return messages, nil
}

func (r *roomResolver) Members(ctx context.Context, obj *model.Room) (*model.Members, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Members", &bson.M{
		"roomID (obj.Rooms)": obj.RoomID,
	})
	defer node.MethodTiming()

	members, err := r.Dataloader.Members(obj.RoomID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных")
	}

	return members, nil
}

// Employee returns generated.EmployeeResolver implementation.
func (r *Resolver) Employee() generated.EmployeeResolver { return &employeeResolver{r} }

// ListenCollection returns generated.ListenCollectionResolver implementation.
func (r *Resolver) ListenCollection() generated.ListenCollectionResolver {
	return &listenCollectionResolver{r}
}

// Me returns generated.MeResolver implementation.
func (r *Resolver) Me() generated.MeResolver { return &meResolver{r} }

// Member returns generated.MemberResolver implementation.
func (r *Resolver) Member() generated.MemberResolver { return &memberResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type employeeResolver struct{ *Resolver }
type listenCollectionResolver struct{ *Resolver }
type meResolver struct{ *Resolver }
type memberResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }
