package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Rooms(ctx context.Context, find model.FindRooms, params *model.Params) (model.RoomsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Rooms", &bson.M{
		"find":   find,
		"params": params,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)

	if node.ValidParams(&params) ||
		find.RoomID != nil && node.ValidID(*find.RoomID) ||
		find.Name != nil && node.ValidRoomName(*find.Name) {
		return node.GetError(), nil
	}

	rooms, err := r.Services.Repos.Rooms.FindRooms(clientID, &find, params)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return rooms, nil
}
