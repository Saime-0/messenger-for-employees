package directive

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/utils"
)

func IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

	if utils.GetAuthDataFromCtx(ctx) == nil {
		err = cerrors.New("не аутентифицирован")
		return obj, err
	}

	return next(ctx)
}

func InputUnion(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	input, ok := obj.(map[string]interface{})
	if !ok {
		panic("InputUnion: can not convert external map")
	}

	valueFound := false
	for _, val := range input {
		if val != nil {
			if valueFound {
				goto handleError
			}
			valueFound = true
		}
	}

	if !valueFound {
		goto handleError
	}

	return next(ctx)

handleError:
	return obj, cerrors.New("one of the input union fields must have a value")
}

func InputLeastOne(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

	input, ok := obj.(map[string]interface{})
	if !ok {
		panic("InputLeastOne: can not convert external map")
	}

	fieldFindIsExists := false
	for key, val := range input {
		if key == "find" || key == "input" {
			fieldFindIsExists = true
			input = val.(map[string]interface{})
			break
		}
	}
	if !fieldFindIsExists {
		panic("InputLeastOne: union input field not found")
	}

	for _, val := range input {
		if fmt.Sprint(val) != "<nil>" {
			return next(ctx)
		}
	}

	return obj, cerrors.New("one of the input fields must have the value")

}
