package middleware

//go:generate mockgen -destination=./mock/idGenerator.go -package=mock -mock_names=IDGenerator=IDGenerator github.com/cloudtrust/common-service/idgenerator IDGenerator
//go:generate mockgen -destination=./mock/logging.go -package=mock -mock_names=Logger=Logger github.com/cloudtrust/common-service/log Logger
//go:generate mockgen -destination=./mock/metrics.go -package=mock -mock_names=Metrics=Metrics,Histogram=Histogram github.com/cloudtrust/common-service/metrics Metrics,Histogram
//go:generate mockgen -destination=./mock/keycloak_client.go -package=mock -mock_names=KeycloakClient=KeycloakClient,IDRetriever=IDRetriever,AdminConfigurationRetriever=AdminConfigurationRetriever github.com/cloudtrust/common-service/middleware KeycloakClient,IDRetriever,AdminConfigurationRetriever
//go:generate mockgen -destination=./mock/tracing.go -package=mock -mock_names=OpentracingClient=OpentracingClient github.com/cloudtrust/common-service/tracing OpentracingClient
