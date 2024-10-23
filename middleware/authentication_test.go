package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cs "github.com/cloudtrust/common-service/v2"
	errorhandler "github.com/cloudtrust/common-service/v2/errors"
	comhttp "github.com/cloudtrust/common-service/v2/http"
	"github.com/cloudtrust/common-service/v2/middleware/mock"
	"github.com/gbrlsnchs/jwt/v2"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	// aud: none
	tokenAudNone = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJJZTVzcXBLdTNwb1g5d1U3YTBhamxnUFlGRHFTTUF5M2l6NEZpelp4d2dnIn0.eyJqdGkiOiJhODg4NTIyNS1kODU5LTRjNDUtODYwZS05YTNjZGYxYjUzZDAiLCJleHAiOjE1NTIyOTQ1NDgsIm5iZiI6MCwiaWF0IjoxNTUyMjkzOTQ4LCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvbWFzdGVyIiwic3ViIjoiNzM5M2FiMWEtNWIwNC00M2Y1LTgwNDktOGE5NDkyMzJlZDBhIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYWRtaW4tY2xpIiwiYXV0aF90aW1lIjowLCJzZXNzaW9uX3N0YXRlIjoiYzdkNTllNTktNTNiYi00Y2IzLThhMTYtZTI3OGI0NWE2OTI5IiwiYWNyIjoiMSIsInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWRtaW4ifQ.WOgsWPdKt1f8gp7AkqCGzoBgkeYgN9YyYlAHILuBG5o9ZN0Ae4Bpymci0tkDWEsQk532mfSyP6-0uLwcNOHf_kPpqjjJ4k6Cnz4p1s6bWTOjPP1cTGcs0bUCiYJI0ZRz3oPjz8RSBH2bDe7Dq7p1STZwLLtX-0uc3t5le0EGSobSoVfOdVBU-TFda4R0xKK7cCsJzw-pOGHFOuoFUhEiruo6Ibo_-iNLxht5rUh8KMoeUkGF3dn1rshT55tq9WY7q6fygUxZS8C_4NvVTfaPo76JO2rUQ5FAhOJRlBACEwALrdpw7Tr0Ox8fjZLIrLeIswMNbGNmpTxEH3LK-ull8g"
	// aud: [] {rpo-realm, test-realm}
	tokenAudArray = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJJZTVzcXBLdTNwb1g5d1U3YTBhamxnUFlGRHFTTUF5M2l6NEZpelp4d2dnIn0.eyJqdGkiOiIwYzYyY2JjMS1hOThlLTQ1NDAtOTM1ZC00NGUwM2M2ZWZkMTAiLCJleHAiOjE1NTY2NjY1MTMsIm5iZiI6MCwiaWF0IjoxNTU2NjMwNTEzLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjpbInJwby1yZWFsbSIsInRlc3QtcmVhbG0iXSwic3ViIjoiNzM5M2FiMWEtNWIwNC00M2Y1LTgwNDktOGE5NDkyMzJlZDBhIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYWRtaW4tY2xpIiwiYXV0aF90aW1lIjowLCJzZXNzaW9uX3N0YXRlIjoiYmEwZjkyNWItZTdhNC00MTBkLWJjY2EtZjU4NzExMWNhOTZlIiwiYWNyIjoiMSIsInJlc291cmNlX2FjY2VzcyI6eyJycG8tcmVhbG0iOnsicm9sZXMiOlsidmlldy1yZWFsbSIsInZpZXctaWRlbnRpdHktcHJvdmlkZXJzIiwibWFuYWdlLWlkZW50aXR5LXByb3ZpZGVycyIsImltcGVyc29uYXRpb24iLCJjcmVhdGUtY2xpZW50IiwibWFuYWdlLXVzZXJzIiwicXVlcnktcmVhbG1zIiwidmlldy1hdXRob3JpemF0aW9uIiwicXVlcnktY2xpZW50cyIsInF1ZXJ5LXVzZXJzIiwibWFuYWdlLWV2ZW50cyIsIm1hbmFnZS1yZWFsbSIsInZpZXctZXZlbnRzIiwidmlldy11c2VycyIsInZpZXctY2xpZW50cyIsIm1hbmFnZS1hdXRob3JpemF0aW9uIiwibWFuYWdlLWNsaWVudHMiLCJxdWVyeS1ncm91cHMiXX0sInRlc3QtcmVhbG0iOnsicm9sZXMiOlsidmlldy1yZWFsbSIsInZpZXctaWRlbnRpdHktcHJvdmlkZXJzIiwibWFuYWdlLWlkZW50aXR5LXByb3ZpZGVycyIsImltcGVyc29uYXRpb24iLCJjcmVhdGUtY2xpZW50IiwibWFuYWdlLXVzZXJzIiwicXVlcnktcmVhbG1zIiwidmlldy1hdXRob3JpemF0aW9uIiwicXVlcnktY2xpZW50cyIsInF1ZXJ5LXVzZXJzIiwibWFuYWdlLWV2ZW50cyIsIm1hbmFnZS1yZWFsbSIsInZpZXctZXZlbnRzIiwidmlldy11c2VycyIsInZpZXctY2xpZW50cyIsIm1hbmFnZS1hdXRob3JpemF0aW9uIiwibWFuYWdlLWNsaWVudHMiLCJxdWVyeS1ncm91cHMiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGdyb3VwcyBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZ3JvdXBzIjpbIi90b2VfYWRtaW5pc3RyYXRvciJdLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImVtYWlsIjoidG90b0B0b3RvLmNvbSJ9.Q62PHuOme8Debm8uhdtvEdMmd5ZX7xrdPfgcgR9MpsInQzykrFZdjUufFQQ1wJw35eaHDdLABXe-IxHPJvqzRS_FrQ54sLGDZz9w6T8umywuSG4VP2UKtJkV7-c1Jswyeq2cbfchteyAsnByXipjXKFYLrWGz5VrtxZKgLbF3lqtLmJzo9RzlEuxbynX0L63kLJism0CWOSxfQuGknMEy9RYp7MmivlHUvisjBMY1lWVyK-cNJUZOyFcANh3PclVrPdZW1QFbynHCnFOfO38vjW7f7Vy2DeGC23YbBG2ZZFRAD7rgM_VfCqjH10w-iGa6G7avOwSD7tGXQCMWLp7Zw"
	// aud: test-realm
	tokenAudString = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJJZTVzcXBLdTNwb1g5d1U3YTBhamxnUFlGRHFTTUF5M2l6NEZpelp4d2dnIn0.eyJqdGkiOiI4MDY4MjZkNy0xZjM4LTQxZjgtYTk5Ni1iYTYzYWI0YTY3MGIiLCJleHAiOjE1NTY2NjY3NzAsIm5iZiI6MCwiaWF0IjoxNTU2NjMwNzcwLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoidGVzdC1yZWFsbSIsInN1YiI6IjczOTNhYjFhLTViMDQtNDNmNS04MDQ5LThhOTQ5MjMyZWQwYSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImFkbWluLWNsaSIsImF1dGhfdGltZSI6MCwic2Vzc2lvbl9zdGF0ZSI6IjFlMmI1Mzk5LTgyNDItNDA1OS05Y2M1LWE5MzI0NDVlY2JkMSIsImFjciI6IjEiLCJyZXNvdXJjZV9hY2Nlc3MiOnsidGVzdC1yZWFsbSI6eyJyb2xlcyI6WyJ2aWV3LXJlYWxtIiwidmlldy1pZGVudGl0eS1wcm92aWRlcnMiLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsImNyZWF0ZS1jbGllbnQiLCJtYW5hZ2UtdXNlcnMiLCJxdWVyeS1yZWFsbXMiLCJ2aWV3LWF1dGhvcml6YXRpb24iLCJxdWVyeS1jbGllbnRzIiwicXVlcnktdXNlcnMiLCJtYW5hZ2UtZXZlbnRzIiwibWFuYWdlLXJlYWxtIiwidmlldy1ldmVudHMiLCJ2aWV3LXVzZXJzIiwidmlldy1jbGllbnRzIiwibWFuYWdlLWF1dGhvcml6YXRpb24iLCJtYW5hZ2UtY2xpZW50cyIsInF1ZXJ5LWdyb3VwcyJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZ3JvdXBzIGVtYWlsIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJncm91cHMiOlsiL3RvZV9hZG1pbmlzdHJhdG9yIl0sInByZWZlcnJlZF91c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJ0b3RvQHRvdG8uY29tIn0.QXUTPciZYYv8k688D27sOz5thyQH1OWwp-rqTnCQYoAbqXPVgSZxLIepk8JvS9drBl7jOH-M_w2tXMOjV-7kY7p57_9VyWaI42VgBVmJVXSWwMwPtWAwnpKqMh1wrrm_zYJRmZ43o1r6Rp_kELnfgwocFSLc3DTDVEoMuYE45kJg9JwPc2K7DYi6Om5qOm9ez-x8GpyGVy3xJiOa-Qr9oJpKCx02sRVEBIc0AE0pfpxfbBhJU06L4uVnwQ1JxquLKLU77bjPEkAKOnTeG-6D9OtH_K42KujZyhj7FytXAXv9CmISi9aIe7BVANFSu7TyOBjelZHVpI5dOKRc-E2L9w"
)

func TestSplitIssuer(t *testing.T) {
	t.Run("Keycloak Wildfly issuer", func(t *testing.T) {
		var issuer, domain = splitIssuer("http://domain/auth/realms/myRealm")
		assert.Equal(t, "http://domain", issuer)
		assert.Equal(t, "myRealm", domain)
	})
	t.Run("Keycloak Quarkus", func(t *testing.T) {
		var issuer, domain = splitIssuer("http://domain/realms/myRealm")
		assert.Equal(t, "http://domain", issuer)
		assert.Equal(t, "myRealm", domain)
	})
}

func TestUnmarshalTokenAudience(t *testing.T) {
	t.Run("Valid token", func(t *testing.T) {
		payload, _, _ := jwt.Parse(tokenAudNone)
		token, err := unmarshalTokenAudience(payload)
		assert.Nil(t, err)
		assert.Equal(t, "admin", token.GetUsername())
	})
	t.Run("Invalid token", func(t *testing.T) {
		_, err := unmarshalTokenAudience([]byte{})
		assert.NotNil(t, err)
	})
}

func getAuthenticationResultTest(m http.Handler, req *http.Request) *http.Response {
	var w = httptest.NewRecorder()
	m.ServeHTTP(w, req)
	return w.Result()
}

func TestHTTPBasicAuthenticationMapMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockLogger = mock.NewLogger(mockCtrl)

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/event/receiver", bytes.NewReader([]byte{}))
	var credentials = map[string]string{
		"jane.doe": "password-doe",
		"john.doe": "password-doe-too",
	}

	var m = MakeHTTPBasicAuthenticationMapMW(credentials, mockLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	t.Run("Missing authorization token", func(t *testing.T) {
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})
	t.Run("Invalid username", func(t *testing.T) {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("whoami:password-doe")))
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})
	t.Run("Invalid password", func(t *testing.T) {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("jane.doe:invalid-password")))
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})
	t.Run("Valid username", func(t *testing.T) {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("jane.doe:password-doe")))
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}

func TestHTTPBasicAuthenticationMW(t *testing.T) {
	var token = "dXNlcm5hbWU6cGFzc3dvcmQ="

	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockLogger = mock.NewLogger(mockCtrl)

	var m = MakeHTTPBasicAuthenticationMW("password", mockLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/event/receiver", bytes.NewReader([]byte{}))

	t.Run("Missing authorization token", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Missing Authorization header").Times(1)
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Non basic format")

	t.Run("Missing basic token", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Missing basic token").Times(1)
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Basic X"+token)
	t.Run("Invalid base64 token", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Invalid base64 token").Times(1)
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Basic "+token)

	t.Run("Valid authorization token - Basic", func(t *testing.T) {
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	req.Header.Set("Authorization", "basic "+token)

	t.Run("Valid authorization token - basic", func(t *testing.T) {
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	req.Header.Set("Authorization", "basic dXNlcm5hbWU6cGFzc3dvcmQx")

	t.Run("Invalid authorization token", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Invalid password value").Times(1)
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Invalid token format (username:password)").Times(1)
		req = httptest.NewRequest("POST", "http://cloudtrust.io/management/test", bytes.NewReader([]byte{}))
		req.Header.Set("Authorization", "Basic 123456ABCDEF")
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Invalid token format (username:password)").Times(1)
		req = httptest.NewRequest("POST", "http://cloudtrust.io/management/test", bytes.NewReader([]byte{}))
		req.Header.Set("Authorization", "Basic dXNlcm5hbWU=")
		var result = getAuthenticationResultTest(m, req)
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})
}

func checkContextEndpoint(ctx context.Context, request any) (response any, err error) {
	var accessToken = ctx.Value(cs.CtContextAccessToken).(string)
	var realm = ctx.Value(cs.CtContextRealm).(string)
	var user = ctx.Value(cs.CtContextUsername).(string)
	var groups = ctx.Value(cs.CtContextGroups).([]string)
	if (tokenAudString == accessToken || tokenAudArray == accessToken) && "master" == realm && "admin" == user && len(groups) == 1 && "toe_administrator" == groups[0] {
		return "", nil
	}

	return "", errorhandler.Error{Status: 500}
}

func TestHTTPOIDCTokenValidationMW(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)

	var m = MakeHTTPOIDCTokenValidationMW(mockKeycloakClient, "test-realm", mockLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/management/test", bytes.NewReader([]byte{}))

	t.Run("Missing authorization token", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Missing Authorization header").Times(1)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Non bearer format")

	t.Run("Missing bearer token", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Missing bearer token").Times(1)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Bearer    AB CD")

	t.Run("Invalid bearer token", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error: Missing bearer token").Times(1)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})

	req.Header.Set("Authorization", "Bearer "+tokenAudString)

	t.Run("Valid authorization token", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockKeycloakClient.EXPECT().VerifyToken(gomock.Any(), "master", tokenAudString).Return(nil).Times(1)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	req.Header.Set("Authorization", "bearer "+tokenAudString)

	t.Run("Invalid authorization token", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error", "err", gomock.Any()).Times(1)
		mockKeycloakClient.EXPECT().VerifyToken(gomock.Any(), "master", tokenAudString).Return(errors.New(errorhandler.MsgErrInvalidParam + "." + errorhandler.Token)).Times(1)
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		var w = httptest.NewRecorder()
		mockLogger.EXPECT().Info(gomock.Any(), "msg", "Authorization error", "err", gomock.Any()).Times(1)
		req = httptest.NewRequest("POST", "http://cloudtrust.io/management/test", bytes.NewReader([]byte{}))
		req.Header.Set("Authorization", "Bearer 123456ABCDEF")
		m.ServeHTTP(w, req)
		var result = w.Result()
		assert.Equal(t, http.StatusForbidden, result.StatusCode)
	})
}

func testAuthentication(t *testing.T, audienceRequired string, token string, expectedStatus int, verifyToken bool) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)

	var handler = http_transport.NewServer(checkContextEndpoint, comhttp.BasicDecodeRequest, comhttp.EncodeReply, http_transport.ServerErrorEncoder(comhttp.ErrorHandlerNoLog()))
	var m = MakeHTTPOIDCTokenValidationMW(mockKeycloakClient, audienceRequired, mockLogger)(handler)

	// HTTP request.
	var req = httptest.NewRequest("POST", "http://cloudtrust.io/management/realms/master", bytes.NewReader([]byte{}))
	req.Header.Set("Authorization", "Bearer "+token)

	var w = httptest.NewRecorder()
	if verifyToken {
		mockKeycloakClient.EXPECT().VerifyToken(gomock.Any(), "master", token).Return(nil).Times(1)
	} else {
		mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	}

	m.ServeHTTP(w, req)
	var result = w.Result()
	assert.Equal(t, expectedStatus, result.StatusCode)
}

func TestContextHTTPOIDCTokenMissingAudience(t *testing.T) {
	testAuthentication(t, "audience", tokenAudNone, http.StatusForbidden, false)
}

func TestContextHTTPOIDCTokenAudienceStringArrayValidationMW(t *testing.T) {
	testAuthentication(t, "rpo-realm", tokenAudArray, http.StatusOK, true)
}

func TestContextHTTPOIDCTokenInvalidAudienceStringArrayMW(t *testing.T) {
	testAuthentication(t, "backoffice", tokenAudArray, http.StatusForbidden, false)
}

func TestContextHTTPOIDCTokenAudienceStringValidationMW(t *testing.T) {
	testAuthentication(t, "test-realm", tokenAudString, http.StatusOK, true)
}

func TestContextHTTPOIDCTokenInvalidAudienceStringMW(t *testing.T) {
	testAuthentication(t, "backoffice", tokenAudString, http.StatusForbidden, false)
}
