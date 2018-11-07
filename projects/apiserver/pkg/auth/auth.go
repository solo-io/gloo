package auth

import (
	"net/http"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/utils/log"
)

var (
	mApiserverGetToken = stats.Int64("apiserver.solo.io/auth/GetToken", "The number of calls to GetToken", "1")
	apiserverGetToken  = &view.View{
		Name:        "apiserver.solo.io/auth/GetToken",
		Measure:     mApiserverGetToken,
		Description: "The number of calls to GetToken",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	mApiserverAuthFail = stats.Int64("apiserver.solo.io/auth/AuthFail", "The number of AuthFails", "1")
	apiserverAuthFail  = &view.View{
		Name:        "apiserver.solo.io/auth/AuthFail",
		Measure:     mApiserverAuthFail,
		Description: "The number of calls to AuthFail",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(apiserverGetToken, apiserverAuthFail)
}

// GetToken returns an oauth bearer token
// The token is stored either in the Authorization header
// or in a cookie
// The header is the primary source of truth and the cookie is
// a cached copy to enable authorization of websocket requests
// (Websockets do not pass the Authorization header)
func GetToken(w http.ResponseWriter, r *http.Request) string {
	stats.Record(r.Context(), mApiserverGetToken.M(1))
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
		stats.Record(r.Context(), mApiserverAuthFail.M(1))
		log.Warnf("could not read bearer token cookie %v\n", err)
		return ""
	}
	return c.Value
}
