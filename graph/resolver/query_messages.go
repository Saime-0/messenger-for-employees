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

func (r *queryResolver) Messages(ctx context.Context, find model.FindMessages, params *model.Params) (model.MessagesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Messages", &bson.M{
		"find": find,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).EmployeeID
	)

	if node.ValidParams(&params) ||
		find.EmpID != nil && node.ValidID(*find.EmpID) ||
		find.RoomID != nil && node.ValidID(*find.RoomID) ||
		node.IsMember(clientID, *find.RoomID) ||
		find.TargetID != nil && node.ValidID(*find.TargetID) ||
		find.TextFragment != nil { // todo text valid
		return node.GetError(), nil
	}

	messages, err := r.Services.Repos.Rooms.FindMessages(clientID, &find, params)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return messages, nil
}
