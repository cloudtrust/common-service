package database

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service Configuration
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/cloudtrustdb.go -package=mock -mock_names=SQLRow=SQLRow,CloudtrustDB=CloudtrustDB,CloudtrustDBFactory=CloudtrustDBFactory github.com/cloudtrust/common-service/database/sqltypes SQLRow,CloudtrustDB,CloudtrustDBFactory
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/database.go -package=mock -mock_names=DbTransactionIntf=DbTransactionIntf github.com/cloudtrust/common-service/database DbTransactionIntf
