package extauth_test_server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	startDiscoveryServerOnce sync.Once
	cachedPrivateKey         *rsa.PrivateKey
	baseOauth2Port           = uint32(5556)
)

type endpointDataFields struct {
	// A map counting the number of times each method has been called
	methodMap map[string]int
}

type endpointData struct {
	data  map[string]*endpointDataFields
	mutex sync.RWMutex
}

func (e *endpointData) incrementMethodCount(path, method string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.data[path] == nil {
		e.data[path] = &endpointDataFields{
			methodMap: make(map[string]int),
		}
	}
	e.getEndpointDataFields(path).methodMap[method] += 1
}

func (e *endpointData) getEndpointDataFields(path string) *endpointDataFields {
	if e.data[path] == nil {
		e.data[path] = &endpointDataFields{
			methodMap: make(map[string]int),
		}
	}
	return e.data[path]
}

// GetMethodCount returns the number of times a given HTTP method was invoked on a given path.
func (e *endpointData) GetMethodCount(path, method string) int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.getEndpointDataFields(path).methodMap[method]
}

type handlerStats struct {
	stats map[string]int
	mutex sync.RWMutex
}

func (h *handlerStats) increment(path string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.stats[path] += 1
}

// Get returns the number of times a path was invoked.
func (h *handlerStats) Get(path string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.stats[path]
}

type FakeDiscoveryServer struct {
	s                        http.Server
	Port                     uint32
	ServerAddress            string
	CreateExpiredToken       bool
	CreateNearlyExpiredToken bool
	Token                    string
	AccessTokenValue         string
	RefreshTokenValue        string
	LastGrant                string
	HandlerStats             *handlerStats
	EndpointData             *endpointData

	// The set of key IDs that are supported by the server
	keyIds []string
}

const (
	AuthEndpoint       = "/auth"
	LogoutEndpoint     = "/logout"
	TokenEndpoint      = "/token"
	RevocationEndpoint = "/revoke"
	KeysEndpoint       = "/keys"
	AlternateAuthPath  = "/alternate-auth"
	ConfigurationPath  = "/.well-known/openid-configuration"
)

func (f *FakeDiscoveryServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = f.s.Shutdown(ctx)
}

func (f *FakeDiscoveryServer) UpdateToken(grantType string) {
	startDiscoveryServerOnce.Do(func() {
		var err error
		cachedPrivateKey, err = rsa.GenerateKey(rand.Reader, 512)
		Expect(err).NotTo(HaveOccurred())
	})

	claims := jwt.MapClaims{
		"foo": "bar",
		"aud": "test-clientid",
		"sub": "user",
		"iss": fmt.Sprintf("http://%s:%d", f.ServerAddress, f.Port),
	}
	f.LastGrant = grantType
	if grantType == "" && f.CreateExpiredToken {
		// create expired token so we can test refresh
		claims["exp"] = time.Now().Add(-time.Minute).Unix()
	} else if grantType == "" && f.CreateNearlyExpiredToken {
		// create token that expires ten ms from now
		claims["exp"] = time.Now().Add(10 * time.Millisecond).Unix()
	}

	tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenToSign.Header["kid"] = f.GetValidKeyId()
	token, err := tokenToSign.SignedString(cachedPrivateKey)
	Expect(err).NotTo(HaveOccurred())
	f.Token = token
}

func (f *FakeDiscoveryServer) UpdateKeyIds(keyIds []string) {
	if len(keyIds) > 0 {
		f.keyIds = keyIds
	}
}

func (f *FakeDiscoveryServer) GetValidKeyId() string {
	// If there is more than one valid kid, return the first one
	return f.keyIds[0]
}

func (f *FakeDiscoveryServer) Start(serverAddress string) *rsa.PrivateKey {
	f.HandlerStats = &handlerStats{
		stats: make(map[string]int),
	}
	f.EndpointData = &endpointData{
		data: make(map[string]*endpointDataFields),
	}
	f.Port = parallel.AdvancePortSafeListen(&baseOauth2Port)
	f.ServerAddress = serverAddress
	// Initialize the server with 1 valid kid
	f.keyIds = []string{"kid-1"}

	f.UpdateToken("")
	n := base64.RawURLEncoding.EncodeToString(cachedPrivateKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(0).SetUint64(uint64(cachedPrivateKey.E)).Bytes())

	f.s = http.Server{
		Addr: fmt.Sprintf(":%d", f.Port),
	}

	f.s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		rw.Header().Set("content-type", "application/json")

		f.HandlerStats.increment(r.URL.Path)
		switch r.URL.Path {
		case AuthEndpoint:
			// redirect back immediately. This simulates a user that's already logged in by the IDP.
			redirectUri := r.URL.Query().Get("redirect_uri")
			state := r.URL.Query().Get("state")
			u, err := url.Parse(redirectUri)
			Expect(err).NotTo(HaveOccurred())

			u.RawQuery = "code=1234&state=" + state
			_, _ = fmt.Fprintf(GinkgoWriter, "redirecting to %s\n", u.String())
			rw.Header().Add("Location", u.String())
			rw.Header().Add("x-auth", "auth")
			rw.WriteHeader(http.StatusFound)

			_, _ = rw.Write([]byte(`auth`))
		case AlternateAuthPath:
			// redirect back immediately. This simulates a user that's already logged in by the IDP.
			redirectUri := r.URL.Query().Get("redirect_uri")
			state := r.URL.Query().Get("state")
			u, err := url.Parse(redirectUri)
			Expect(err).NotTo(HaveOccurred())

			u.RawQuery = "code=9876&state=" + state
			_, _ = fmt.Fprintf(GinkgoWriter, "redirecting to %s\n", u.String())
			rw.Header().Add("Location", u.String())
			rw.Header().Add("x-auth", "alternate-auth")
			rw.WriteHeader(http.StatusFound)

			_, _ = rw.Write([]byte(`alternate-auth`))
		case ConfigurationPath:
			appUrl := fmt.Sprintf("http://%s:%d", f.ServerAddress, f.Port)
			authUrl := fmt.Sprintf("%s%s", appUrl, AuthEndpoint)
			tokenUrl := fmt.Sprintf("%s%s", appUrl, TokenEndpoint)
			revocationUrl := fmt.Sprintf("%s%s", appUrl, RevocationEndpoint)
			logoutUrl := fmt.Sprintf("%s%s", appUrl, LogoutEndpoint)
			keysUrl := fmt.Sprintf("%s%s", appUrl, KeysEndpoint)
			_, _ = rw.Write([]byte(fmt.Sprintf(`
		{
			"issuer": "%s",
			"authorization_endpoint": "%s",
			"token_endpoint": "%s",
			"revocation_endpoint": "%s",
			"end_session_endpoint": "%s",
			"jwks_uri": "%s",
			"response_types_supported": [
			  "code"
			],
			"subject_types_supported": [
			  "public"
			],
			"id_token_signing_alg_values_supported": [
			  "RS256"
			],
			"scopes_supported": [
			  "openid",
			  "email",
			  "profile"
			]
		  }
		`, appUrl, authUrl, tokenUrl, revocationUrl, logoutUrl, keysUrl)))
		case TokenEndpoint:
			err := r.ParseForm()
			Expect(err).NotTo(HaveOccurred())
			_, _ = fmt.Fprintln(GinkgoWriter, "got request for token. query:", r.URL.RawQuery, r.URL.String(), "form:", r.Form.Encode(), "\n full request", r)
			if r.URL.Query().Get("grant_type") == "refresh_token" || r.Form.Get("grant_type") == "refresh_token" {
				f.UpdateToken("refresh_token")
			}
			_, _ = rw.Write([]byte(`
			{
				"access_token": "` + f.AccessTokenValue + `",
				"token_type": "Bearer",
				"refresh_token": "` + f.RefreshTokenValue + `",
				"expires_in": 3600,
				"id_token": "` + f.Token + `"
			 }
	`))
		case KeysEndpoint:
			var keyListBuffer bytes.Buffer
			for _, kid := range f.keyIds {
				keyListBuffer.WriteString(`
				{
					"use": "sig",
					"kty": "RSA",
					"kid": "` + kid + `",
					"alg": "RS256",
					"n": "` + n + `",
					"e": "` + e + `"
				},`)
			}
			// Remove the last comma so it's valid json
			keyList := strings.TrimSuffix(keyListBuffer.String(), ",")
			keysResponse := `
			{
				"keys": [
				    ` + keyList + `
				]
			}
			`
			_, _ = rw.Write([]byte(keysResponse))
		case RevocationEndpoint:
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
			_, _ = rw.Write([]byte(httpReply))
		case LogoutEndpoint: // this should match the end_session_endpoint path noted
			// in the /.well-known/openid-configuration above.
			// This should be invoked by the Ext-Auth service when the client
			// calls the logoutPath on the OAuth Configuration, and has a session.
			// A GET request will add the client_id and the id_token_hint onto the query parameters.
			// A POST request will add them to the request's Form.
			var clientId, token string
			switch r.Method {
			case http.MethodGet:
				queryParams := r.URL.Query()
				clientId = queryParams.Get("client_id")
				token = queryParams.Get("id_token_hint")
			case http.MethodPost:
				defer r.Body.Close()
				body, err := io.ReadAll(r.Body)
				Expect(err).NotTo(HaveOccurred())

				// The response is in the form of `client_id=...&id_token_hint=...`, so we need to parse it to extract those fields.
				bodyStr := string(body)
				bodyData := strings.Split(bodyStr, "&")

				bodyDataMap := make(map[string]string)
				for _, data := range bodyData {
					dataParts := strings.Split(data, "=")
					bodyDataMap[dataParts[0]] = dataParts[1]
				}

				clientId = bodyDataMap["client_id"]
				token = bodyDataMap["id_token_hint"]
			}

			// "test-clientid"
			Expect(len(clientId)).Should(BeNumerically(">", 0))
			Expect(len(token)).Should(BeNumerically(">", 0))

			f.EndpointData.incrementMethodCount(LogoutEndpoint, r.Method)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(""))
		}
	})

	go func() {
		defer GinkgoRecover()
		err := f.s.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}()

	return cachedPrivateKey
}

func (f *FakeDiscoveryServer) GetOidcAuthCodeConfig(envoyPort uint32, appUrl string, secretRef *core.ResourceRef) *extauth.OAuth2_OidcAuthorizationCode {
	return &extauth.OAuth2_OidcAuthorizationCode{
		OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
			ClientId:        "test-clientid",
			ClientSecretRef: secretRef,
			IssuerUrl:       fmt.Sprintf("http://%s:%d/", f.ServerAddress, f.Port),
			AppUrl:          fmt.Sprintf("http://%s:%d", appUrl, envoyPort),
			CallbackPath:    "/callback",
			LogoutPath:      LogoutEndpoint,
			Scopes:          []string{"email"},
		},
	}
}

func (f *FakeDiscoveryServer) GetOauthConfig(envoyPort uint32, appUrl string, secretRef *core.ResourceRef) *extauth.OAuth {
	return &extauth.OAuth{
		ClientId:        "test-clientid",
		ClientSecretRef: secretRef,
		IssuerUrl:       fmt.Sprintf("http://%s:%d/", f.ServerAddress, f.Port),
		AppUrl:          fmt.Sprintf("http://%s:%d", appUrl, envoyPort),
		CallbackPath:    "/callback",
		Scopes:          []string{"email"},
	}
}

// Helper function that generates an id token using the key id passed in and a private key. The issuer is automatically
// set to the fake discovery server's address and port which is initialized when starting the server.
func (f *FakeDiscoveryServer) GenerateIdTokenWithKid(kid string, privateKey *rsa.PrivateKey) string {
	tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"foo": "bar",
		"aud": "test-clientid",
		"sub": "user",
		"iss": fmt.Sprintf("http://%s:%d", f.ServerAddress, f.Port),
	})
	tokenToSign.Header["kid"] = kid
	token, err := tokenToSign.SignedString(privateKey)
	Expect(err).NotTo(HaveOccurred())

	return token
}

func (f *FakeDiscoveryServer) GenerateValidIdToken(privateKey *rsa.PrivateKey) string {
	return f.GenerateIdTokenWithKid(f.GetValidKeyId(), privateKey)
}
