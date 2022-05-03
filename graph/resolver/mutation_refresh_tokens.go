package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/models"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/internal/rules"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) RefreshTokens(ctx context.Context, sessionKey *string, refreshToken string) (model.RefreshTokensResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("RefreshTokens", &bson.M{
		"sessionKey":   sessionKey,
		"refreshToken": refreshToken,
	})
	defer node.MethodTiming()

	sessionID, clientID, err := r.Services.Repos.Auth.FindSessionByComparedToken(refreshToken)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}
	if sessionID == 0 {
		return resp.Error(resp.ErrBadRequest, "неверный токен"), nil
	}
	var (
		session *models.RefreshSession
	)
	newRefreshToken := kit.RandomSecret(rules.RefreshTokenLength)
	sessionExpAt := kit.After(*r.Config.RefreshTokenLifetime)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(res.CtxUserAgent).(string),
		ExpAt:        sessionExpAt,
	}

	err = r.Services.Repos.Auth.UpdateRefreshSession(sessionID, session)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	tokenExpiresAt := kit.After(*r.Config.AccessTokenLifetime)
	token, err := utils.GenerateToken(
		&utils.TokenData{
			EmployeeID: clientID,
			ExpiresAt:  tokenExpiresAt,
		},
		r.Config.SecretSigningKey,
	)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	if sessionKey != nil {
		err = r.Subix.ExtendClientSession(*sessionKey, tokenExpiresAt)
		if err != nil {
			node.Healer.Debug(cerrors.Wrap(err, utils.GetCallerPos()+""))
		}
	}

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); ok && sessionExpAt < runAt.(int64) {
		_, err = r.Services.Scheduler.AddTask(
			func() {
				err := r.Services.Repos.Employees.DeleteRefreshSession(sessionID)
				if err != nil {

					r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
				}
			},
			sessionExpAt,
		)
		if err != nil {
			panic(err)
		}
	}

	return model.TokenPair{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}
