package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/piper"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) SetNotify(ctx context.Context, roomID int, value bool) (model.SetNotifyResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("SetNotify", &bson.M{
		"roomID": roomID,
		"value":  value,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)

	if node.RoomExists(roomID) ||
		node.IsMember(clientID, roomID) {
		return node.GetError(), nil
	}

	err := func(n piper.Node) error {
		n.SwitchMethod("CreateMessage", &bson.M{
			"clientID": clientID,
			"roomID":   roomID,
			"value":    value,
		})
		defer n.MethodTiming()

		return r.Services.Repos.Rooms.SetNotify(clientID, roomID, value)
	}(node)

	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить настройку"), nil
	}

	return resp.Success("успешно"), nil
}
