package directives

import (
	"context"
	"myapp/middleware"
	"myapp/service"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func IsLogin(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	var (
		user = middleware.AuthContext(ctx)
		s    = service.GetService()
	)

	if user == nil {
		return nil, &gqlerror.Error{
			Message: "User not logged in",
			Extensions: map[string]interface{}{
				"code": http.StatusUnauthorized,
			},
		}
	}

	exist, _ := s.UserCheckExistByID(ctx, user.ID)
	if !exist {
		return nil, &gqlerror.Error{
			Message: "User not found",
			Extensions: map[string]interface{}{
				"code": http.StatusNotFound,
			},
		}
	}

	return next(ctx)
}
