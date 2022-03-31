package utils

import (
	"context"
	"github.com/robbert229/jwt"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/res"
)

type TokenData struct {
	EmployeeID int   `json:"employeeid"`
	ExpiresAt  int64 `json:"exp"`
}

func ParseToken(tokenString string, secretKey string) (*TokenData, error) {

	var (
		employeeID  int
		expiresAt   int64
		data        *TokenData
		err         error
		claims      *jwt.Claims
		_employeeID interface{}
		_expiresAt  interface{}
		femployeeID float64
		fexpiresAt  float64
		ok          bool
		algorithm   jwt.Algorithm
	)

	algorithm = jwt.HmacSha256(secretKey)
	if err := algorithm.Validate(tokenString); err != nil {
		goto handleError
	}

	claims, err = algorithm.Decode(tokenString)
	if err != nil {
		goto handleError
	}

	_employeeID, err = claims.Get("employeeid")
	if err != nil {
		goto handleError
	}
	_expiresAt, err = claims.Get("exp")
	if err != nil {
		goto handleError
	}

	femployeeID, ok = _employeeID.(float64)
	if !ok {
		err = cerrors.New("token not contain employeeid")
		goto handleError
	}
	fexpiresAt, ok = _expiresAt.(float64)
	if !ok {
		err = cerrors.New("token not contain exp")
		goto handleError
	}

	employeeID = int(femployeeID)
	expiresAt = int64(fexpiresAt)

	data = &TokenData{
		EmployeeID: employeeID,
		ExpiresAt:  expiresAt,
	}

	return data, nil

handleError:
	return nil, cerrors.Wrap(err, "не удалось распарсить токен")
}

func GenerateToken(data *TokenData, secretKey string) (string, error) {
	algorithm := jwt.HmacSha256(secretKey)

	claims := jwt.NewClaim()
	claims.Set("employeeid", data.EmployeeID)
	claims.Set("exp", data.ExpiresAt)

	token, err := algorithm.Encode(claims)

	if err != nil {
		return "", cerrors.Wrap(err, "не удалось сгенерировать токен")
	}

	return token, nil
}

func GetAuthDataFromCtx(ctx context.Context) (authData *TokenData) {
	data, ok := ctx.Value(res.CtxAuthData).(*TokenData)
	if !ok {
		println("GetAuthDataFromCtx: не удалось найти CtxAuthData в контексте")
	}
	return data
}
