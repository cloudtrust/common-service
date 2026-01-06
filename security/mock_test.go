package security

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient,IdentificationKeycloakClient=IdentificationKeycloakClient github.com/cloudtrust/common-service/v2/security KeycloakClient,IdentificationKeycloakClient
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/authentication_db_reader.go -package=mock -mock_names=AuthorizationDBReader=AuthorizationDBReader,IdentificationAuthorizationDBReader=IdentificationAuthorizationDBReader github.com/cloudtrust/common-service/v2/security AuthorizationDBReader,IdentificationAuthorizationDBReader
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/detailederr.go -package=mock -mock_names=DetailedError=DetailedError github.com/cloudtrust/common-service/v2/errors DetailedError
