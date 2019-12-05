package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"
	"strings"

	cs "github.com/cloudtrust/common-service"
	errorhandler "github.com/cloudtrust/common-service/errors"
	"github.com/cloudtrust/common-service/log"
	"github.com/gbrlsnchs/jwt"
)

// MakeHTTPBasicAuthenticationMW retrieve the token from the HTTP header 'Basic' and
// check if the password value match the allowed one.
// If there is no such header, the request is not allowed.
// If the password is correct, the username is added into the context:
//   - username: username extracted from the token
func MakeHTTPBasicAuthenticationMW(passwordToMatch string, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var authorizationHeader = req.Header.Get("Authorization")
			var ctx = context.TODO()

			if authorizationHeader == "" {
				logger.Info(ctx, "Authorization Error", "Missing Authorization header")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.AuthHeader), w)
				return
			}

			var regexpBasicAuth = `^[Bb]asic (.+)$`
			var r = regexp.MustCompile(regexpBasicAuth)
			var match = r.FindStringSubmatch(authorizationHeader)
			if match == nil {
				logger.Info(ctx, "Authorization Error", "Missing basic token")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.BasicToken), w)
				return
			}

			// Decode base 64
			decodedToken, err := base64.StdEncoding.DecodeString(match[1])

			if err != nil {
				logger.Info(ctx, "Authorization Error", "Invalid base64 token")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}

			// Extract username & password values
			var tokenSubparts = strings.Split(string(decodedToken), ":")

			if len(tokenSubparts) != 2 {
				logger.Info(ctx, "Authorization Error", "Invalid token format (username:password)")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}

			var username = tokenSubparts[0]
			var password = tokenSubparts[1]

			ctx = context.WithValue(req.Context(), cs.CtContextUsername, username)

			// Check password match
			if password != passwordToMatch {
				logger.Info(ctx, "Authorization Error", "Invalid password value")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// KeycloakClient is the interface of the keycloak client.
type KeycloakClient interface {
	VerifyToken(ctx context.Context, realmName string, accessToken string) error
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
				logger.Info(ctx, "Authorization Error", "Missing Authorization header")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.AuthHeader), w)
				return
			}

			var r = regexp.MustCompile(`^[Bb]earer +([^ ]+)$`)
			var match = r.FindStringSubmatch(authorizationHeader)
			if match == nil {
				logger.Info(ctx, "Authorization Error", "Missing bearer token")
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrMissingParam+"."+errorhandler.BearerToken), w)
				return
			}

			// match[0] is the global matched group. match[1] is the first captured group
			var accessToken = match[1]

			payload, _, err := jwt.Parse(accessToken)
			if err != nil {
				logger.Info(ctx, "Authorization Error", err)
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}

			var userID, username, issuer, issuerDomain, realm string
			var groups []string

			// The audience in JWT may be a string array or a string.
			// First we try with a string array, if a failure occurs we try with a string
			{
				var jot TokenAudienceStringArray
				if err = jwt.Unmarshal(payload, &jot); err == nil {
					userID = jot.Subject
					username = jot.Username
					issuer = jot.Issuer
					var splitIssuer = strings.Split(issuer, "/auth/realms/")
					issuerDomain = splitIssuer[0]
					realm = splitIssuer[1]
					groups = extractGroups(jot.Groups)

					if !assertMatchingAudience(jot.Audience, audienceRequired) {
						logger.Info(ctx, "Authorization Error", "Incorrect audience")
						httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
						return
					}
				}
			}

			if err != nil {
				var jot TokenAudienceString
				if err = jwt.Unmarshal(payload, &jot); err == nil {
					userID = jot.Subject
					username = jot.Username
					issuer = jot.Issuer
					var splitIssuer = strings.Split(issuer, "/auth/realms/")
					issuerDomain = splitIssuer[0]
					realm = splitIssuer[1]
					groups = extractGroups(jot.Groups)

					if jot.Audience != audienceRequired {
						logger.Info(ctx, "Authorization Error", "Incorrect audience")
						httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
						return
					}
				} else {
					logger.Info(ctx, "Authorization Error", err)
					httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
					return
				}
			}

			ctx = context.WithValue(req.Context(), cs.CtContextAccessToken, accessToken)
			ctx = context.WithValue(ctx, cs.CtContextRealm, realm)
			ctx = context.WithValue(ctx, cs.CtContextUserID, userID)
			ctx = context.WithValue(ctx, cs.CtContextUsername, username)
			ctx = context.WithValue(ctx, cs.CtContextGroups, groups)
			ctx = context.WithValue(ctx, cs.CtContextIssuerDomain, issuerDomain)

			if err = keycloakClient.VerifyToken(ctx, realm, accessToken); err != nil {
				logger.Info(ctx, "Authorization Error", err)
				httpErrorHandler(ctx, http.StatusForbidden, errors.New(errorhandler.MsgErrInvalidParam+"."+errorhandler.Token), w)
				return
			}

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func assertMatchingAudience(jwtAudiences []string, requiredAudience string) bool {
	for _, jwtAudience := range jwtAudiences {
		if requiredAudience == jwtAudience {
			return true
		}
	}

	return false
}

func extractGroups(kcGroups []string) []string {
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

type header struct {
	Algorithm   string `json:"alg,omitempty"`
	KeyID       string `json:"kid,omitempty"`
	Type        string `json:"typ,omitempty"`
	ContentType string `json:"cty,omitempty"`
}
