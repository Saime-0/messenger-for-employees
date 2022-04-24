package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) RoomMessages(ctx context.Context, roomID int, startMsg int, created model.MsgCreated, count int) (model.MessagesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("RoomMessages", &bson.M{
		"roomID":   roomID,
		"startMsg": startMsg,
		"created":  created,
		"count":    count,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)

	if node.ValidID(roomID) ||
		node.IsMember(clientID, roomID) ||
		node.ValidID(startMsg) ||
		node.ValidMsgCount(count) {
		return node.GetError(), nil
	}

	messages, err := r.Services.Repos.Rooms.RoomMessages(roomID, startMsg, created, count)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return messages, nil
}
