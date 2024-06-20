package middleware

import (
	"context"
	"fmt"
	"joshsoftware/peerly/config"
	"joshsoftware/peerly/service"
	"net/http"
	"strconv"
	"strings"

	ae "joshsoftware/peerly/apperrors"

	jwt "github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
)


func JwtAuthMiddleware(next http.Handler, deps service.Dependencies) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add orgid, userid, roleid to request context
		// ctx := context.WithValue(r.Context(), "orgid", orgID)
		ctx := context.WithValue(r.Context(), "orgid", 1)
		// ctx = context.WithValue(ctx, "userid", userID)
		ctx = context.WithValue(ctx, "userid", 1)
		// ctx = context.WithValue(ctx, "roleid", roleID)
		ctx = context.WithValue(ctx, "roleid", 1)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
		// next.ServeHTTP(w, r)
		return
		totalFields := 2
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != totalFields {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}
		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.JWTKey()), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx = context.WithValue(r.Context(), "props", claims)
		userID, err := strconv.Atoi(claims["sub"].(string))
		if err != nil {
			logger.Error(ae.ErrJSONParseFail, "Error parsing JSON for token response", err)
			return
		}

		orgID, err := strconv.Atoi(claims["org"].(string))
		if err != nil {
			logger.Error(ae.ErrJSONParseFail, "Error parsing JSON for token response", err)
			return
		}

		currentUser, err := deps.Store.GetUser(ctx, userID)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while fetching User")
			return
		}

		nextContext := context.WithValue(ctx, "user", token)
		if currentUser.OrgID != orgID {
			err = ae.ErrInvalidToken
			logger.WithField("err", err.Error()).Error("Mismatch with user organization and current organization")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(nextContext))
	})
}

func GetUserId () int64{
	return 0
}

func GetRole() int64{
	return 1
}