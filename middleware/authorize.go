package middleware

import (
	"context"
	"errors"
	"fmt"
	"myginapp/common"
	"myginapp/component/tokenprovider"
	"myginapp/modules/user/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthenStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"),
	)
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	// "Authorization": "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}

	return parts[1], nil
}

func RequiredAuth(authStore AuthenStore, tokenProvider tokenprovider.Provider) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			panic(err)
		}

		// db := appCtx.GetMa inDBConnection()
		// store := userstore.NewSQLStore(db)

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
			// handle error
		}
		// user, err := store.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId})

		user, err := authStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId()})
		if err != nil {
			panic(err)
			// handle error
		}

		if user.Status == 0 {
			panic(common.ErrNoPermission(errors.New("user has been deleted or banned")))
		}

		c.Set(common.CurrentUser, user)
		c.Next()
	}
}
