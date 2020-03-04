package configuration

//go:generate mockgen -destination=./mock/cloudtrustdb.go -package=mock -mock_names=CloudtrustDB=CloudtrustDB,SQLRow=SQLRow github.com/cloudtrust/common-service/database/sqltypes CloudtrustDB,SQLRow
