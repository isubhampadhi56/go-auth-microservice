package authMiddleware

import (
	"context"
	"net/http"
	"strings"

	tokencache "github.com/go-auth-microservice/pkg/model/tokenCache"
	jwtauth "github.com/go-auth-microservice/pkg/utils/jwtAuth"
	"github.com/go-auth-microservice/pkg/utils/logger"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

func AccessTokenVerify(next http.Handler) http.Handler {
	var blackListedToken tokencache.BlackListedToken = tokencache.GetBlacklistTokenCache()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.InitializeAuditLogger()
		accessTokenHandler := jwtauth.GetAccessTokenHandler()
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" || !strings.HasPrefix(accessToken, "Bearer ") {
			http.Error(w, "missing or invalid access token", http.StatusUnauthorized)
			log.Error("access token not present or invalid token")
			return
		}
		if blackListedToken.IsPresent(accessToken) {
			http.Error(w, "token expired or user account has been updated", http.StatusUnauthorized)
			log.Error("token expired or user account has been updated")
			return
		}
		claims, err := accessTokenHandler.VerifyToken(accessToken)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			log.Error(err)
			return
		}
		userId, _ := claims["userId"].(float64)
		ctx := context.WithValue(r.Context(), contextKey("userId"), uint64(userId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
