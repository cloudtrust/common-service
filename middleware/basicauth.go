package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/cloudtrust/common-service/v2/log"
)

const (
	regexpBasicAuthHeader = `(?i)^basic\s+([\w+\/]+=*)\s*$`
)

var (
	// ErrBasicNotSupported is returned when the authorization header is not supported
	ErrBasicNotSupported = errors.New("not supported")
	// ErrBasicUnknownUser is returned when basic authorization is received but user/password is invalid
	ErrBasicUnknownUser = errors.New("unknown user")
)

// BasicAuthCollection struct
type BasicAuthCollection struct {
	basicAuths []string
}

// NewBasicAuthCollection creates a new instance of BasicAuthCollection
func NewBasicAuthCollection() *BasicAuthCollection {
	return &BasicAuthCollection{}
}

func (b *BasicAuthCollection) getBasicAuthFromHeader(authorization string) (string, error) {
	re := regexp.MustCompile(regexpBasicAuthHeader)
	matches := re.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return "", ErrBasicNotSupported
	}
	return matches[1], nil
}

func (b *BasicAuthCollection) getUserFromBasicAuth(basicAuth string) string {
	decoded, err := base64.StdEncoding.DecodeString(basicAuth)
	if err != nil {
		return "(invalid base64)" // User will be associated to an error
	}
	parts := strings.Split(string(decoded), ":")
	if len(parts) < 2 {
		return "(invalid basic auth format)" // User will be associated to an error
	}
	return parts[0]
}

// IsEmpty checks if the BasicAuthCollection is empty
func (b *BasicAuthCollection) IsEmpty() bool {
	return len(b.basicAuths) == 0
}

// Import imports a JSON string containing basic auth users and passwords into the BasicAuthCollection
func (b *BasicAuthCollection) Import(jsonContent string) error {
	var basicAuthUsersMap map[string]string
	err := json.Unmarshal([]byte(jsonContent), &basicAuthUsersMap)
	if err != nil {
		return errors.New("allowed users cannot be parsed")
	}
	for pseudo, password := range basicAuthUsersMap {
		b.Add(pseudo, password)
	}
	return nil
}

// ImportFromMap imports users from a map of user/password
func (b *BasicAuthCollection) ImportFromMap(creds map[string]string) {
	for user, password := range creds {
		b.Add(user, password)
	}
}

// Add adds a new user/password to the BasicAuthCollection
func (b *BasicAuthCollection) Add(pseudo, password string) {
	secret := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", pseudo, password))
	b.basicAuths = append(b.basicAuths, strings.ReplaceAll(secret, "=", ""))
}

// AllowedUserFromAuthorizationHeader gets allowed user from the HTTP authorization header
func (b *BasicAuthCollection) AllowedUserFromAuthorizationHeader(authorization string) (*string, error) {
	basicAuth, err := b.getBasicAuthFromHeader(authorization)
	if err != nil {
		return nil, err
	}
	return b.AllowedUserFromAuthorizationValue(basicAuth)
}

// AllowedUserFromAuthorizationValue gets allowed user from the value of the "Basic" content of the authorization header
func (b *BasicAuthCollection) AllowedUserFromAuthorizationValue(basicAuth string) (*string, error) {
	user := b.getUserFromBasicAuth(basicAuth)
	if !slices.Contains(b.basicAuths, strings.ReplaceAll(basicAuth, "=", "")) {
		return nil, ErrBasicUnknownUser
	}
	return &user, nil
}

// MakeHTTPBasicAuthMW creates a middleware to retrieve the token from the HTTP header 'Basic' and check credentials
func (b *BasicAuthCollection) MakeHTTPBasicAuthMW(logger log.Logger) func(http.Handler) http.Handler {
	return MakeHTTPBasicAuthenticationFuncMW(b.AllowedUserFromAuthorizationValue, logger)
}
