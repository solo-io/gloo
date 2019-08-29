package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"

	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/go-ldap/ldap"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/errors"
	"go.uber.org/zap"
)

var (
	// see: https://ldapwiki.com/wiki/Best%20Practices%20For%20LDAP%20Naming%20Attributes#section-Best+Practices+For+LDAP+Naming+Attributes-SpecialCharacters
	specialLdapCharsRegex = regexp.MustCompile("[,=+<>#;\"]")
	FailedConnectionError = func(err error) error {
		return errors.Wrapf(err, "failed to connect to LDAP service")
	}
	FailedSearchError = func(err error) error {
		return errors.Wrapf(err, "failed to search LDAP server")
	}
)

type ldapAuthService struct {
	// Used to build connections to the LDAP server
	clientBuilder ClientBuilder
	// LDAP server address
	serverUrl string
	// Template to build user entry distinguished names (DN), e.g. "uid=%s,ou=people,dc=solo,dc=io"
	userDnTemplate string
	// User must be member of one of these groups for the request to be authenticated
	allowedGroups map[string]bool
}

func NewLdapAuthService(clientBuilder ClientBuilder, config *Config) *ldapAuthService {

	// Index groups for fast access
	allowedGroups := make(map[string]bool)
	for _, g := range config.AllowedGroups {
		allowedGroups[g] = true
	}

	return &ldapAuthService{
		clientBuilder:  clientBuilder,
		serverUrl:      config.ServerUrl,
		userDnTemplate: config.UserDnTemplate,
		allowedGroups:  allowedGroups,
	}
}

func (s *ldapAuthService) Start(ctx context.Context) error {
	conn, err := s.clientBuilder.Dial("tcp", s.serverUrl)
	if err != nil {
		return FailedConnectionError(err)
	}
	conn.Close()
	return nil
}

func (s *ldapAuthService) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {

	// Get required header values
	username, password, ok := GetBasicAuthCredentials(request)
	if !ok {
		ur := unauthenticatedResponse()
		logger(ctx).Infow("failed to retrieve basic auth credentials from header", zap.Any("request", request), zap.Any("response", ur))
		return ur, nil
	}

	// Remove special characters to protect against potential injection attacks
	sanitizedUsername, wasSanitized := SanitizeLdapDN(username)
	if wasSanitized {
		logger(ctx).Warnw("username contained special characters and has been sanitized",
			zap.Any("original", username),
			zap.Any("sanitized", sanitizedUsername),
		)
	}

	// TODO(marco): use connection pool
	// Establish connection to LDAP server
	conn, err := s.clientBuilder.Dial("tcp", s.serverUrl)
	if err != nil {
		return nil, FailedConnectionError(err)
	}
	defer conn.Close()

	// Build user entry distinguished name
	userDN := fmt.Sprintf(s.userDnTemplate, sanitizedUsername)

	// Perform BIND. If it fails, we consider authentication to have failed and return a 401 response.
	//
	// An LDAP session (connection from a client) has an authorization state. When a session is established,
	// the authorization state is set to anonymous. LDAP clients use the BIND operation to change the authorization
	// state of a session.
	// When the directory server receives a BIND request from a client, the authorization state of that connection is
	// set to anonymous. If the BIND request is successful, the authorization state of the connection is set to the
	// state for the identity in the BIND request; if the BIND request is not successful, the session remains in the
	// anonymous state. An LDAP directory server considers a BIND operation to be successful if the information
	// transmitted in the BIND request can be verified by the server.
	if err = conn.Bind(userDN, password); err != nil {
		ur := unauthenticatedResponse()
		logger(ctx).Infow("failed to BIND to LDAP server", zap.Any("userDN", userDN), zap.Any("response", ur))
		return ur, nil
	}

	// Try to retrieve the 'memberOf' attribute for the user's own entry
	searchReq := ldap.NewSearchRequest(
		userDN,               // The base dn to search
		ldap.ScopeBaseObject, // Use 'base' scope (get only the base DN entry)
		ldap.NeverDerefAliases,
		0, 0, false,
		"(uid=*)",                     // The filter to apply, needs to not be empty
		[]string{MembershipAttribute}, // Attributes to retrieve
		nil,
	)

	// TODO(marco): turn infos to debugs
	result, err := conn.Search(searchReq)
	if err != nil {
		logger(ctx).Infow("failed to search LDAP server", zap.Any("userDN", userDN), zap.Any("searchRequest", searchReq))
		return nil, FailedSearchError(err)
	}

	for _, entry := range result.Entries {
		memberOf := entry.GetAttributeValues(MembershipAttribute)
		for _, group := range memberOf {
			if _, ok := s.allowedGroups[group]; ok {
				logger(ctx).Infow("user is a member of allowed group",
					zap.Any("userDN", userDN),
					zap.Any("allowedGroups", s.allowedGroups),
					zap.Any("userMemberships", memberOf),
					zap.Any("matchedGroup", group),
				)
				return api.AuthorizedResponse(), nil
			}
		}
	}

	return api.UnauthorizedResponse(), nil
}

func SanitizeLdapDN(input string) (result string, wasSanitized bool) {
	if specialLdapCharsRegex.Match([]byte(input)) {
		sanitized := specialLdapCharsRegex.ReplaceAllString(input, "\\${0}")
		return sanitized, true
	}
	return input, false
}

func GetBasicAuthCredentials(request *api.AuthorizationRequest) (string, string, bool) {
	headers := request.CheckRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	if headers == nil {
		return "", "", false
	}

	// Header is in the form: "authorization: Basic <credentials>"
	authHeaderValue, ok := headers["authorization"]
	if !ok {
		// try upper case before failing
		authHeaderValue, ok = headers["Authorization"]
		if !ok {
			return "", "", false
		}
	}

	return parseBasicAuth(authHeaderValue)
}

// Copied over from http.Request
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func unauthenticatedResponse() *api.AuthorizationResponse {
	return &api.AuthorizationResponse{
		CheckResponse: envoyauthv2.CheckResponse{
			Status: &rpc.Status{
				Code: int32(rpc.UNAUTHENTICATED),
			},
			HttpResponse: &envoyauthv2.CheckResponse_DeniedResponse{
				DeniedResponse: &envoyauthv2.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_Unauthorized,
					},
				},
			},
		},
	}
}
