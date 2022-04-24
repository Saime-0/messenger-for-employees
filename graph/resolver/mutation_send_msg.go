package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/models"
	"github.com/saime-0/messenger-for-employee/internal/piper"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) SendMsg(ctx context.Context, input model.CreateMessageInput) (model.SendMsgResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("SendMsg", &bson.M{
		"Rooms":       input.RoomID,
		"TargetMsgID": input.TargetMsgID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)

	if node.RoomExists(input.RoomID) ||
		node.IsMember(clientID, input.RoomID) ||
		input.TargetMsgID != nil && node.MessageExists(input.RoomID, *input.TargetMsgID) {
		return node.GetError(), nil
	}

	message := &models.CreateMessage{
		TargetMsgID: input.TargetMsgID,
		EmployeeID:  clientID,
		RoomID:      input.RoomID,
		Body:        input.Body,
	}

	eventReadyMessage, err := func(n piper.Node) (*model.NewMessage, error) {
		n.SwitchMethod("CreateMessage", &bson.M{
			"TargetMsgID": message.TargetMsgID,
			"EmployeeID":  message.EmployeeID,
			"Rooms":       message.RoomID,
		})
		defer n.MethodTiming()

		return r.Services.Repos.Rooms.CreateMessage(message)
	}(node)

	if err != nil {

		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "не удалось создать сообщение"), nil
	}

	//r.Services.Events.NewMessage(roomID, &model.Message{Rooms:      msgID, TargetMsgID: _replyTo, Rooms:  &model.Member{Rooms: memberID}, Type:    message.Type, Body:    input.Body})
	go func() {
		err := r.Subix.NotifyRoomMembers(
			eventReadyMessage,
			input.RoomID,
		)
		if err != nil {

			node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		}
	}()

	return resp.Success("сообщение успешно создано"), nil
}
