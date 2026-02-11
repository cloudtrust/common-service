package middleware

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/v2/idgenerator IDGenerator
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/cloudtrust/common-service/v2/log Logger
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient,IDRetriever=IDRetriever,AdminConfigurationRetriever=AdminConfigurationRetriever github.com/cloudtrust/common-service/v2/middleware KeycloakClient,IDRetriever,AdminConfigurationRetriever
