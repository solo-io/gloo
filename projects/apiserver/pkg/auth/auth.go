package auth

import (
	"net/http"

	"github.com/solo-io/solo-kit/pkg/utils/log"
)

// GetToken returns an oauth bearer token
// The token is stored either in the Authorization header
// or in a cookie
// The header is the primary source of truth and the cookie is
// a cached copy to enable authorization of websocket requests
// (Websockets do not pass the Authorization header)
func GetToken(w http.ResponseWriter, r *http.Request) string {
	newToken := r.Header.Get("Authorization")
	if newToken != "" {
		setBearerCookie(w, r, newToken)
		return newToken
	}
	return getBearerCookie(r)
}

func setBearerCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "cBearer",
		Value: token,
		// Secure: true, // TODO(mitchdraft) - use secure cookies when we switch to https
	})
}

func getBearerCookie(r *http.Request) string {
	c, err := r.Cookie("cBearer")
	if err != nil {
		log.Warnf("could not read bearer token cookie %v\n", err)
		return ""
	}
	return c.Value
}
