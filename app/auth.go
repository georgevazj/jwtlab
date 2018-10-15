package app

import (
	"net/http"
	"github.com/georgevazj/jwtlab/utils"
	"strings"
	"github.com/georgevazj/jwtlab/models"
	"github.com/dgrijalva/jwt-go"
	"os"
	"fmt"
	"context"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// List of endpoints that doesn't require auth
		notAuth := []string{"/api/user/new", "/api/user/login"}

		// Current request path
		requestPath := r.URL.Path

		// Check if request does not need authentication , serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		// Grab the token from the header
		tokenHeader := r.Header.Get("Authorization")
		// If token is missing returns with error code 403 Unauthorized
		if tokenHeader == "" {
			response = utils.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		// The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requiremen
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type","application/json")
			utils.Respond(w, response)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error){
			return []byte(os.Getenv("token_password")), nil
		})

		// Malformed token, returns with http code 403 as usual
		if err != nil {
			response = utils.Message(false, "Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		// Token is invalid, maybe not signed on this server
		if !token.Valid {
			response = utils.Message(false, "Token is not valid")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		// Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %d", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}