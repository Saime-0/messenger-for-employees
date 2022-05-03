package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *subscriptionResolver) Subscribe(ctx context.Context, sessionKey string) (<-chan *model.SubscriptionBody, error) {

	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Subscribe", &bson.M{
		"sessionKey": sessionKey,
	})
	defer node.MethodTiming()

	node.Debug(&bson.M{
		"sessionKey": sessionKey,
	})

	var (
		authData = utils.GetAuthDataFromCtx(ctx)
	)

	if authData == nil { // тк @isAuth  вебсокетинге не отрабатывает
		node.Debug("не аутентифицирован")
		return nil, cerrors.New("не аутентифицирован")
	}

	if node.ValidSessionKey(sessionKey) {
		node.Debug(cerrors.Wrap(cerrors.New(node.GetError().Error), utils.GetCallerPos()+""))
		return nil, cerrors.New(node.GetError().Error)
	}

	client, err := r.Subix.Sub(
		authData.EmployeeID,
		sessionKey,
		authData.ExpiresAt,
	)
	if err != nil {
		node.Debug(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return nil, err
	}

	// New client
	go func() {
		<-ctx.Done()
		// client is down
		r.Subix.Unsub(sessionKey)
		log.Printf("удаляю клиента") // debug
	}()

	return client.Ch, nil
}
