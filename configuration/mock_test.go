package configuration

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/cloudtrustdb.go -package=mock -mock_names=CloudtrustDB=CloudtrustDB,SQLRow=SQLRow,SQLRows=SQLRows github.com/cloudtrust/common-service/database/sqltypes CloudtrustDB,SQLRow,SQLRows
