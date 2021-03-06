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

func (r *queryResolver) Tags(ctx context.Context, tagIDs []int, params *model.Params) (model.TagsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Tags", &bson.M{
		"params": params,
	})
	defer node.MethodTiming()

	var (
	//clientID = utils.GetAuthDataFromCtx(ctx).Rooms
	)

	if node.ValidParams(&params) {
		return node.GetError(), nil
	}

	tags, err := r.Services.Repos.Tags.Tags(tagIDs, params)
	if err != nil {

		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return tags, nil
}
