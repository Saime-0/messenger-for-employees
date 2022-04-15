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

func (r *mutationResolver) MoveRoom(ctx context.Context, roomID int, prevRoomID *int) (model.MoveRoomResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("MoveRoom", &bson.M{
		"roomID":     roomID,
		"prevRoomID": prevRoomID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)
	if node.RoomExists(roomID) ||
		prevRoomID != nil && node.RoomExists(*prevRoomID) ||
		node.IsMember(clientID, roomID) ||
		prevRoomID != nil && node.IsMember(clientID, *prevRoomID) {
		return node.GetError(), nil
	}

	err := r.Services.Repos.Rooms.MoveRoom(clientID, roomID, prevRoomID)

	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось переместить комнату"), nil
	}

	return resp.Success("комната успешно перемешена"), nil
}
