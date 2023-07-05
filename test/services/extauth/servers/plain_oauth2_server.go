package extauth_test_server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	expectedIdsAndSecrets = map[string]string{"test-clientid": "test", "no-secret-id": ""}
)

type FakeOAuth2Server struct {
	s                 http.Server
	Token             string
	RefreshTokenValue string
	AccessTokenValue  string
	Port              uint32
}

func (f *FakeOAuth2Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = f.s.Shutdown(ctx)
}

func (f *FakeOAuth2Server) Start() {
	f.Port = atomic.AddUint32(&baseOauth2Port, 1) + uint32(parallel.GetPortOffset())
	f.s = http.Server{
		Addr: fmt.Sprintf(":%d", f.Port),
	}

	f.s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		rw.Header().Set("content-type", "text/plain")

		switch r.URL.Path {
		case "/auth":
			// redirect back immediately. This simulates a user that's already logged in by the IDP.
			redirectUri := r.URL.Query().Get("redirect_uri")
			state := r.URL.Query().Get("state")
			glooUrlRedirect := r.URL.Query().Get("gloo_urlToRedirect")
			u, err := url.Parse(redirectUri)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			u.RawQuery = fmt.Sprintf("code=1234&state=%s&gloo_urlToRedirect=%s", state, glooUrlRedirect)
			_, err = fmt.Fprintf(GinkgoWriter, "redirecting to %s\n", u.String())
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			rw.Header().Add("Location", u.String())
			rw.Header().Add("x-auth", "auth")
			rw.WriteHeader(http.StatusFound)

			_, _ = rw.Write([]byte(`auth`))
		case "/token":
			err := r.ParseForm()
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			clientId, clientSecret, _ := r.BasicAuth()
			expectedSecret, ok := expectedIdsAndSecrets[clientId]
			if !ok || expectedSecret != clientSecret {
				rw.WriteHeader(http.StatusUnauthorized)
			} else {
				f.Token = f.AccessTokenValue
				refreshToken := f.RefreshTokenValue

				_, err = rw.Write([]byte(fmt.Sprintf("access_token=%s&refresh_token=%s", f.Token, refreshToken)))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
			}
		case "/revoke":
			err := r.ParseForm()
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			tokenTypeHint := r.Form.Get("token_type_hint")

			httpReply := ""
			if tokenTypeHint != "refresh_token" && tokenTypeHint != "access_token" {
				httpReply = `
              {
              	"error":"unsupported_token_type"
				}
              `
			}

			rw.WriteHeader(http.StatusOK)
			_, err = rw.Write([]byte(httpReply))
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		}
	})

	go func() {
		defer GinkgoRecover()
		err := f.s.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}()
}

func (f *FakeOAuth2Server) GetOAuth2Config(envoyPort uint32, url string, secretRef *core.ResourceRef) *extauth.OAuth2_Oauth2 {
	return &extauth.OAuth2_Oauth2{
		Oauth2: &extauth.PlainOAuth2{
			ClientId:        "test-clientid",
			ClientSecretRef: secretRef,
			AppUrl:          fmt.Sprintf("http://%s:%d", url, envoyPort),
			CallbackPath:    "/callback",
			Scopes:          []string{"email"},
			AuthEndpoint:    fmt.Sprintf("http://%s:%d/auth", url, f.Port),
			TokenEndpoint:   fmt.Sprintf("http://%s:%d/token", url, f.Port),
		},
	}
}
