package events

import (
	"context"
	"io"
	"log"

	"github.com/IBM/sarama"
	cloudtrust_log "github.com/cloudtrust/common-service/v2/log"
)

func NewSaramaLogger(logger cloudtrust_log.Logger, enabled bool) sarama.StdLogger {
	if enabled {
		return log.New(&cloudtrustLoggerWrapper{logger}, "[Sarama] ", log.LstdFlags)
	}
	return log.New(io.Discard, "[Sarama] ", log.LstdFlags)
}

type cloudtrustLoggerWrapper struct {
	logger cloudtrust_log.Logger
}

func (c *cloudtrustLoggerWrapper) Write(p []byte) (n int, err error) {
	c.logger.Info(context.Background(), "saramaMsg", string(p))
	return len(p), nil
}
