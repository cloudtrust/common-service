package healthcheck

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/healthcheck.go -package=mock -mock_names=HealthDatabase=HealthDatabase github.com/cloudtrust/common-service/healthcheck HealthDatabase
