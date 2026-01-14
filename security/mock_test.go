package security

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient,RoleBasedKeycloakClient=RoleBasedKeycloakClient github.com/cloudtrust/common-service/v2/security KeycloakClient,RoleBasedKeycloakClient
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/authentication_db_reader.go -package=mock -mock_names=AuthorizationDBReader=AuthorizationDBReader,RoleBasedAuthorizationDBReader=RoleBasedAuthorizationDBReader github.com/cloudtrust/common-service/v2/security AuthorizationDBReader,RoleBasedAuthorizationDBReader
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/detailederr.go -package=mock -mock_names=DetailedError=DetailedError github.com/cloudtrust/common-service/v2/errors DetailedError
