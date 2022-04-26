package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/rules"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) RoomMessages(ctx context.Context, byCreated *model.ByCreated, byRange *model.ByRange) (model.MessagesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("RoomMessagesByCreated", &bson.M{
		"byCreated": byCreated,
		"byRange":   byRange,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
		err      error
		messages *model.Messages
	)

	if byCreated != nil && (node.ValidID(byCreated.RoomID) ||
		node.IsMember(clientID, byCreated.RoomID) ||
		node.ValidID(byCreated.StartMsg) ||
		node.ValidMsgCount(byCreated.Count)) ||
		byRange != nil && (node.ValidID(byRange.RoomID) ||
			node.IsMember(clientID, byRange.RoomID) ||
			node.ValidID(byRange.Start) || node.ValidID(byRange.End)) {
		return node.GetError(), nil
	}

	if byCreated != nil {
		messages, err = r.Services.Repos.Rooms.RoomMessagesByCreated(byCreated)
	} else if byRange != nil {
		messages, err = r.Services.Repos.Rooms.RoomMessagesByRange(byRange, rules.MaxMsgCount)
	} else {
		return resp.Error(resp.ErrBadRequest, "требуется ввести хотя бы один параметр"), nil
	}

	if err != nil {

		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return messages, nil
}
