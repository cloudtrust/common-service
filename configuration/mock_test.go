package configuration

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/cloudtrustdb.go -package=mock -mock_names=CloudtrustDB=CloudtrustDB,SQLRow=SQLRow,SQLRows=SQLRows,Transaction github.com/cloudtrust/common-service/v2/database/sqltypes CloudtrustDB,SQLRow,SQLRows
