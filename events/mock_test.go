package events

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/sarama.go -package=mock -mock_names=SyncProducer=SyncProducer github.com/IBM/sarama SyncProducer
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/log.go -package=mock -mock_names=Logger=Logger github.com/cloudtrust/common-service/v2/log Logger
