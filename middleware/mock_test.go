package middleware

//generate mockgen --build_flags=--mod=mod -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/idgenerator IDGenerator
//generate mockgen --build_flags=--mod=mod -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/cloudtrust/common-service/log Logger
//generate mockgen --build_flags=--mod=mod -destination=./mock/metrics.go -package=mock -mock_names=Metrics=Metrics,Histogram=Histogram github.com/cloudtrust/common-service/metrics Metrics,Histogram
//generate mockgen --build_flags=--mod=mod -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient,IDRetriever=IDRetriever,AdminConfigurationRetriever=AdminConfigurationRetriever github.com/cloudtrust/common-service/middleware KeycloakClient,IDRetriever,AdminConfigurationRetriever
//generate mockgen --build_flags=--mod=mod -destination=./mock/tracing.go -package=mock -mock_names=OpentracingClient=OpentracingClient github.com/cloudtrust/common-service/tracing OpentracingClient
