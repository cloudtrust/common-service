package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	cs "github.com/cloudtrust/common-service"
	errorhandler "github.com/cloudtrust/common-service/errors"
	"github.com/cloudtrust/common-service/log"
	"github.com/cloudtrust/common-service/security"
	"github.com/gbrlsnchs/jwt/v2"
	errorsPkg "github.com/pkg/errors"
)

// MakeHTTPBasicAuthenticationFuncMW retrieve the token from the HTTP header 'Basic' and
// check credentials according to the given callback function
// If there is no such header, the request is not allowed.
// If the password is correct, the username is added into the context
func MakeHTTPBasicAuthenticationFuncMW(credsMatcher func(token string) (*string, error), logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var ctx = context.TODO()
			var token, err = extractBasicAuthentication(ctx, req.Header.Get("Authorization"), logger)
			if err != nil {
				httpErrorHandler(ctx, http.StatusForbidden, err, w)
				return
			}

			var authenticated *string
			if authenticated, err = credsMatcher(token); err != nil {
				httpErrorHandler(ctx, http.StatusForbidden, err, w)
				return
			} else if authenticated == nil {
				logger.Info(ctx, "msg", "Authorization error: Invalid password value")
				httpErrorHandler(ctx, http.StatusUnauthorized, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}
			ctx = context.WithValue(req.Context(), cs.CtContextUsername, *authenticated)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// MakeHTTPBasicAuthenticationMapMW retrieve the token from the HTTP header 'Basic' and
// check credentials according to the given credentials map
// If there is no such header, the request is not allowed.
// If the password is correct, the username is added into the context
func MakeHTTPBasicAuthenticationMapMW(credentials map[string]string, logger log.Logger) func(http.Handler) http.Handler {
	var authTokens = make(map[string]string)
	for user, password := range credentials {
		var token = fmt.Sprintf("%s:%s", user, password)
		var token64 = base64.StdEncoding.EncodeToString([]byte(token))
		authTokens[token64] = user
	}

	return MakeHTTPBasicAuthenticationFuncMW(func(token string) (*string, error) {
		if username, ok := authTokens[token]; ok {
			return &username, nil
		}
		return nil, nil
	}, logger)
}

// MakeHTTPBasicAuthenticationMW retrieve the token from the HTTP header 'Basic' and
// check if the password value match the allowed one.
// If there is no such header, the request is not allowed.
// If the password is correct, the username is added into the context:
//   - username: username extracted from the token
func MakeHTTPBasicAuthenticationMW(passwordToMatch string, logger log.Logger) func(http.Handler) http.Handler {
	return MakeHTTPBasicAuthenticationFuncMW(func(token string) (*string, error) {
		var ctx = context.TODO()
		var username, password, err = decodeBasicAuthToken(ctx, token, logger)
		if err != nil {
			return nil, err
		}
		if password == passwordToMatch {
			return &username, nil
		}
		return nil, nil
	}, logger)
}

func extractBasicAuthentication(ctx context.Context, authorizationHeader string, logger log.Logger) (string, error) {
	if authorizationHeader == "" {
		logger.Info(ctx, "msg", "Authorization error: Missing Authorization header")
		return "", errors.New(errorhandler.MsgErrMissingParam + "." + errorhandler.AuthHeader)
	}

	var regexpBasicAuth = `^[Bb]asic (.+)$`
	var r = regexp.MustCompile(regexpBasicAuth)
	var match = r.FindStringSubmatch(authorizationHeader)
	if match == nil {
		logger.Info(ctx, "msg", "Authorization error: Missing basic token")
		return "", errors.New(errorhandler.MsgErrMissingParam + "." + errorhandler.BasicToken)
	}

	return match[1], nil
}

func decodeBasicAuthToken(ctx context.Context, authToken string, logger log.Logger) (string, string, error) {
	// Decode base 64
	decodedToken, err := base64.StdEncoding.DecodeString(authToken)

	if err != nil {
		logger.Info(ctx, "msg", "Authorization error: Invalid base64 token")
		return "", "", errors.New(errorhandler.MsgErrInvalidParam + "." + errorhandler.Token)
	}

	// Extract username & password values
	var tokenSubparts = strings.Split(string(decodedToken), ":")

	if len(tokenSubparts) != 2 {
		logger.Info(ctx, "msg", "Authorization error: Invalid token format (username:password)")
		return "", "", errors.New(errorhandler.MsgErrInvalidParam + "." + errorhandler.Token)
	}

	return tokenSubparts[0], tokenSubparts[1], nil
}

// KeycloakClient is the interface of the keycloak client.
type KeycloakClient interface {
	VerifyToken(issuer string, realmName string, accessToken string) error
}

// MakeHTTPOIDCTokenValidationMW retrieve the oidc token from the HTTP header 'Bearer' and
// check its validity for the Keycloak instance binded to the component.
// If there is no such header, the request is not allowed.
// If the token is validated, the following informations are added into the context:
//   - access_token: the recieved access token in raw format
//   - realm: realm name extracted from the Issuer information of the token
//   - username: username extracted from the token
func MakeHTTPOIDCTokenValidationMW(keycloakClient KeycloakClient, audienceRequired string, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var authorizationHeader = req.Header.Get("Authorization")
			var ctx = context.TODO()

			if authorizationHeader == "" {
				logger.Info(ctx, "msg", "Authorization error: Missing Authorization header")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.AuthHeader), w)
				return
			}

			var r = regexp.MustCompile(`^[Bb]earer +([^ ]+)$`)
			var match = r.FindStringSubmatch(authorizationHeader)
			if match == nil {
				logger.Info(ctx, "msg", "Authorization error: Missing bearer token")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.BearerToken), w)
				return
			}

			// match[0] is the global matched group. match[1] is the first captured group
			var accessToken = match[1]

			var jot TokenAudience

			jot, err := ParseAndValidateOIDCToken(ctx, accessToken, keycloakClient, audienceRequired, logger)

			// If there was an error during the validation process, raise an error and stop
			if err != nil {
				switch errorsPkg.Cause(err).(type) {
				case security.ForbiddenError:
					httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
					break
				case errorhandler.UnauthorizedError:
					httpErrorHandler(ctx, http.StatusUnauthorized, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
					break
				}
				return
			}

			var issuer, issuerDomain, realm string
			issuer = jot.GetIssuer()
			var splitIssuer = strings.Split(issuer, "/auth/realms/")
			issuerDomain = splitIssuer[0]
			realm = splitIssuer[1]

			ctx = context.WithValue(req.Context(), cs.CtContextAccessToken, accessToken)
			ctx = context.WithValue(ctx, cs.CtContextRealm, realm)
			ctx = context.WithValue(ctx, cs.CtContextUserID, jot.GetSubject())
			ctx = context.WithValue(ctx, cs.CtContextUsername, jot.GetUsername())
			ctx = context.WithValue(ctx, cs.CtContextGroups, ExtractGroups(jot.GetGroups()))
			ctx = context.WithValue(ctx, cs.CtContextIssuerDomain, issuerDomain)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// ParseAndValidateOIDCToken ensures the OIDC token given in parameter is valid. This method must be public as it is used externally by some projects
func ParseAndValidateOIDCToken(ctx context.Context, accessToken string, keycloakClient KeycloakClient, audienceRequired string, logger log.Logger) (TokenAudience, error) {

	payload, _, err := jwt.Parse(accessToken)
	if err != nil {
		logger.Info(ctx, "msg", "Authorization error", "err", err)
		return nil, security.ForbiddenError{}
	}

	var jot TokenAudience

	if jot, err = unmarshalTokenAudience(payload); err != nil {
		logger.Info(ctx, "msg", "Authorization error", "err", err)
		return nil, security.ForbiddenError{}
	}

	if !jot.AssertMatchingAudience(audienceRequired) {
		logger.Info(ctx, "msg", "Authorization error: Incorrect audience")
		return nil, security.ForbiddenError{}
	}

	var issuer, issuerDomain, realm string
	issuer = jot.GetIssuer()
	var splitIssuer = strings.Split(issuer, "/auth/realms/")
	issuerDomain = splitIssuer[0]
	realm = splitIssuer[1]

	if err = keycloakClient.VerifyToken(issuerDomain, realm, accessToken); err != nil {
		logger.Info(ctx, "msg", "Authorization error", "err", err)
		return nil, errorhandler.UnauthorizedError{}
	}

	// if there was no error during the token validation process, return true
	return jot, nil
}

// AssertMatchingAudience checks if the required audience is in the jwt list of audiences
func AssertMatchingAudience(jwtAudiences []string, requiredAudience string) bool {
	for _, jwtAudience := range jwtAudiences {
		if requiredAudience == jwtAudience {
			return true
		}
	}

	return false
}

// ExtractGroups extracts the list of groups
func ExtractGroups(kcGroups []string) []string {
	var groups = []string{}

	for _, kcGroup := range kcGroups {
		groups = append(groups, strings.TrimPrefix(kcGroup, "/"))
	}

	return groups
}

// TokenAudienceStringArray is JWT token and the custom fields present in OIDC Token provided by Keycloak.
// Audience can be a string or a string array according the specification.
// The libraries are not supporting tit at this time (Fix in progress), meanwhile we circumvent it with a quick fix.
type TokenAudienceStringArray struct {
	hdr            *header
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       []string `json:"aud,omitempty"`
	ExpirationTime int64    `json:"exp,omitempty"`
	NotBefore      int64    `json:"nbf,omitempty"`
	IssuedAt       int64    `json:"iat,omitempty"`
	ID             string   `json:"jti,omitempty"`
	Username       string   `json:"preferred_username,omitempty"`
	Groups         []string `json:"groups,omitempty"`
}

// TokenAudienceString is JWT token with an Audience field represented as a string
type TokenAudienceString struct {
	hdr            *header
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       string   `json:"aud,omitempty"`
	ExpirationTime int64    `json:"exp,omitempty"`
	NotBefore      int64    `json:"nbf,omitempty"`
	IssuedAt       int64    `json:"iat,omitempty"`
	ID             string   `json:"jti,omitempty"`
	Username       string   `json:"preferred_username,omitempty"`
	Groups         []string `json:"groups,omitempty"`
}

type TokenAudience interface {
	GetSubject() string
	GetUsername() string
	GetIssuer() string
	GetGroups() []string

	AssertMatchingAudience(requiredValue string) bool
}

type header struct {
	Algorithm   string `json:"alg,omitempty"`
	KeyID       string `json:"kid,omitempty"`
	Type        string `json:"typ,omitempty"`
	ContentType string `json:"cty,omitempty"`
}

func unmarshalTokenAudience(payload []byte) (TokenAudience, error) {
	var err error

	// The audience in JWT may be a string array or a string.
	// First we try with a string array, if a failure occurs we try with a string
	{
		var jot TokenAudienceStringArray
		if err = jwt.Unmarshal(payload, &jot); err == nil {
			return &jot, nil
		}
	}

	{
		var jot TokenAudienceString
		if err = jwt.Unmarshal(payload, &jot); err == nil {
			return &jot, nil
		}
	}
	return nil, err
}

// GetSubject provides the subject from the token
func (ta *TokenAudienceStringArray) GetSubject() string { return ta.Subject }

// GetUsername provides the username from the token
func (ta *TokenAudienceStringArray) GetUsername() string { return ta.Username }

// GetIssuer provides the issuer from the token
func (ta *TokenAudienceStringArray) GetIssuer() string { return ta.Issuer }

// GetGroups provides the groups from the token
func (ta *TokenAudienceStringArray) GetGroups() []string { return ta.Groups }

// AssertMatchingAudience checks if the required audience is in the token list of audiences
func (ta *TokenAudienceStringArray) AssertMatchingAudience(requiredValue string) bool {
	return AssertMatchingAudience(ta.Audience, requiredValue)
}

// GetSubject provides the subject from the token
func (ta *TokenAudienceString) GetSubject() string { return ta.Subject }

// GetUsername provides the username from the token
func (ta *TokenAudienceString) GetUsername() string { return ta.Username }

// GetIssuer provides the issuer from the token
func (ta *TokenAudienceString) GetIssuer() string { return ta.Issuer }

// GetGroups provides the groups from the token
func (ta *TokenAudienceString) GetGroups() []string { return ta.Groups }

// AssertMatchingAudience checks if the required audience is in the token list of audiences
func (ta *TokenAudienceString) AssertMatchingAudience(requiredValue string) bool {
	return ta.Audience == requiredValue
}
