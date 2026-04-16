package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const CurrentUserKey contextKey = "currentUser"

type AuthUser struct {
	ID       string
	Username string
	Role     string
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			sub, _ := claims["sub"].(string)
			username, _ := claims["username"].(string)
			role, _ := claims["role"].(string)

			user := &AuthUser{
				ID:       sub,
				Username: username,
				Role:     role,
			}

			ctx := context.WithValue(r.Context(), CurrentUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCurrentUserFromCTX(ctx context.Context) (*AuthUser, error) {
	if ctx.Value(CurrentUserKey) == nil {
		return nil, fmt.Errorf("no user in context")
	}

	user, ok := ctx.Value(CurrentUserKey).(*AuthUser)
	if !ok || user.ID == "" {
		return nil, fmt.Errorf("no user in context")
	}

	return user, nil
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header != "" && len(header) > 7 && strings.ToUpper(header[:7]) == "BEARER " {
		return header[7:]
	}

	if token := r.URL.Query().Get("access_token"); token != "" {
		return token
	}

	return ""
}
