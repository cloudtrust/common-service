package events

import (
	"context"

	"github.com/cloudtrust/common-service/v2/log"

	"github.com/Shopify/sarama"
	cs "github.com/cloudtrust/common-service/v2"
)

type KafkaProducerConfig struct {
	Version      string
	Brokers      []string
	Topic        string
	ClientID     string
	ClientSecret string
	TokenURL     string
	Noop         bool
}

func GetKafkaProducerConfig(c cs.Configuration, prefix string) KafkaProducerConfig {
	var cfg KafkaProducerConfig

	cfg.Noop = !c.GetBool(prefix)

	if !cfg.Noop {
		cfg.Version = c.GetString(prefix + "-version")
		cfg.Brokers = c.GetStringSlice(prefix + "-brokers")
		cfg.Topic = c.GetString(prefix + "-topic")
		cfg.ClientID = c.GetString(prefix + "-client-id")
		cfg.ClientSecret = c.GetString(prefix + "-client-secret")
		cfg.TokenURL = c.GetString(prefix + "-token-url")
	}

	return cfg
}

func NewEventKafkaProducer(ctx context.Context, c KafkaProducerConfig, logger log.Logger) (sarama.SyncProducer, error) {
	if c.Noop {
		return &NoopKafkaProducer{}, nil
	}

	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		logger.Error(ctx, "msg", "Error parsing Kafka version", "error", err)
		return nil, err
	}
	config := sarama.NewConfig()
	config.Version = version

	config.Producer.Return.Successes = true

	// Enables Oauth2 authentification
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	config.Net.SASL.TokenProvider = NewTokenProvider(c.ClientID, c.ClientSecret, c.TokenURL)

	producer, err := sarama.NewSyncProducer(c.Brokers, config)
	if err != nil {
		logger.Error(ctx, "msg", "Failed to start Kafka producer", "error", err)
		return nil, err
	}
	return producer, nil
}
