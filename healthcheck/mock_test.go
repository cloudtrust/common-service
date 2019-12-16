package healthcheck

//go:generate mockgen -destination=./mock/healthcheck.go -package=mock -mock_names=HealthDatabase=HealthDatabase github.com/cloudtrust/common-service/healthcheck HealthDatabase
