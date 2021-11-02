package passage

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// AuthenticateRequest determines whether or not to authenticate via header or cookie authentication
// returns the userID (string) on success, error on failure
func (a *App) AuthenticateRequest(r *http.Request) (string, error) {
	if a.Config.HeaderAuth {
		return a.AuthenticateRequestWithHeader(r)
	}
	return a.AuthenticateRequestWithCookie(r)
}

// AuthenticateRequestWithCookie fetches the bearer token from the authorization header and uses it to authenticate
// returns the userID (string) on success, error on failure
func (a *App) AuthenticateRequestWithHeader(r *http.Request) (string, error) {
	authHeaderFields := strings.Fields(r.Header.Get("Authorization"))
	if len(authHeaderFields) != 2 || authHeaderFields[0] == "Bearer" {
		return "", errors.New("missing authentication token: expected \"Bearer\" header")
	}

	userID, valid := a.ValidateAuthToken(authHeaderFields[1])
	if !valid {
		return "", errors.New("invalid authentication token")
	}

	return userID, nil
}

// AuthenticateRequestWithCookie fetches a cookie from the request and uses it to authenticate
// returns the userID (string) on success, error on failure
func (a *App) AuthenticateRequestWithCookie(r *http.Request) (string, error) {
	authTokenCookie, err := r.Cookie("psg_auth_token")
	if err != nil {
		return "", errors.New("missing authentication token: expected \"psg_auth_token\" cookie")
	}

	userID, valid := a.ValidateAuthToken(authTokenCookie.Value)
	if !valid {
		return "", errors.New("invalid authentication token")
	}

	return userID, nil
}

// ValidateAuthToken determines whether a JWT is valid or not
// returns userID (string) on success, error on failure
func (a *App) ValidateAuthToken(authToken string) (string, bool) {
	// Verify that the authentication token is valid:
	parsedToken, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing algorithm")
		}
		return a.PublicKey, nil
	})
	if err != nil {
		return "", false
	}

	// Extract claims from JWT:
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", false
	}

	return userID, true
}
