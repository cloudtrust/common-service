package healthcheck

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/healthcheck.go -package=mock -mock_names=HealthDatabase=HealthDatabase github.com/cloudtrust/common-service/v2/healthcheck HealthDatabase
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/eventsreportermodule.go -package=mock -mock_names=AuditEventsReporterModule=AuditEventsReporterModule github.com/cloudtrust/common-service/v2/events AuditEventsReporterModule
