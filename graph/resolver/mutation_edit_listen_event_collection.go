package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) EditListenEventCollection(ctx context.Context, sessionKey string, action model.EventSubjectAction, targetRooms []int, listenEvents []model.EventType) (model.EditListenEventCollectionResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("EditListenEventCollection", &bson.M{
		"sessionKey":   sessionKey,
		"action":       action,
		"targetRooms":  targetRooms,
		"listenEvents": listenEvents,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).EmployeeID

	if len(listenEvents) > len(model.AllEventType) {
		return resp.Error(resp.ErrBadRequest, "недопустимая длина списка событий"), nil
	}
	targetRooms = kit.GetUniqueInts(targetRooms) // избавляемся от повторяющихся значений

	if node.ValidSessionKey(sessionKey) ||
		node.EmployeeHasAccessToRooms(clientID, targetRooms) {
		return node.GetError(), nil
	}

	err := r.Subix.ModifyCollection(clientID, sessionKey, targetRooms, action, listenEvents)
	if err != nil {
		return resp.Error(resp.ErrBadRequest, "не удалось обновить коллекцию"), nil
	}

	return &model.ListenCollection{SessionKey: sessionKey}, nil
}
