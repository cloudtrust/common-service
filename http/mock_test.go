package http

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/kit_logger.go -package=mock -mock_names=Logger=Logger "github.com/go-kit/kit/log" Logger
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/responsewriter.go -package=mock -mock_names=ResponseWriter=ResponseWriter net/http ResponseWriter
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/detailederr.go -package=mock -mock_names=DetailedError=DetailedError github.com/cloudtrust/common-service/v2/errors DetailedError
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/authorization.go -package=mock -mock_names=AuthorizationManager=AuthorizationManager github.com/cloudtrust/common-service/v2/security AuthorizationManager
