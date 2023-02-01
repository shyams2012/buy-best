package interfaces

import (
	"context"
	"errors"

	"log"
	"net/http"
	"strings"

	"github.com/shyams2012/buy-best/graph/model"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func MiddlewareGetUserFromToken(next http.Handler, resolver *Resolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		//validate jwt token
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))

		user, err := model.ParseAuthToken(token)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			log.Printf("[ERROR] %s", err.Error())
			return
		}

		var userFromDB model.User
		if tx := resolver.DB().Where("id = ?", user.ID).First(&userFromDB); tx.Error != nil {
			http.Error(w, "Invalid user token", http.StatusForbidden)
			log.Printf("[ERROR] %s", tx.Error)
			return
		}
		// put it in context
		ctx := context.WithValue(r.Context(), UserCtxKey, user)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// UserForContext finds the user from the context. REQUIRES MiddlewareGetUserFromToken to have run.
func UserForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(UserCtxKey).(*model.User)
	return raw
}

func CheckAuth(ctx context.Context, role []model.UserRole) (*model.User, error) {

	users := UserForContext(ctx)
	if users == nil {
		return nil, errors.New("need authentication")
	}
	for _, v := range role {
		if users.Role == v {
			return users, nil
		}
	}
	return users, nil
}
