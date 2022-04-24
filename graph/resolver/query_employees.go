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

func (r *queryResolver) Employees(ctx context.Context, find model.FindEmployees, params *model.Params) (model.EmployeesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Employees", &bson.M{
		"find":   find,
		"params": params,
	})
	defer node.MethodTiming()

	if node.ValidParams(&params) ||
		find.EmpID != nil && node.ValidID(*find.EmpID) ||
		find.RoomID != nil && node.ValidID(*find.RoomID) ||
		find.TagID != nil && node.ValidID(*find.TagID) ||
		find.Name != nil && node.ValidNameFragment(*find.Name) {
		return node.GetError(), nil
	}

	employees, err := r.Services.Repos.Employees.FindEmployees(&find)
	if err != nil {

		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return employees, nil
}
