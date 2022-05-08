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

func (r *mutationResolver) ReadMsg(ctx context.Context, roomID int, msgID int) (model.ReadMsgResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("ReadMsg", &bson.M{
		"roomID": roomID,
		"msgID":  msgID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)
	if node.RoomExists(roomID) ||
		node.IsMember(clientID, roomID) ||
		node.MessageExists(roomID, msgID) {
		return node.GetError(), nil
	}

	err := r.Services.Repos.Rooms.ReadMessage(clientID, roomID, msgID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "не удалось прочитать сообщение"), nil
	}

	return resp.Success("успех"), nil
}
