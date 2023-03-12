package gkafka

import (
	"context"
	"github.com/Shopify/sarama"
)

type Options struct {
	Hosts                   []string
	Topics                  []string
	GroupID                 string
	Retry                   bool
	Handle                  func(ctx context.Context, message *sarama.ConsumerMessage) error
	ProducerSuccessCallback func(message *sarama.ProducerMessage)
	ProducerErrorCallback   func(message *sarama.ProducerError)
	ConsumerErrorCallback   func(message error)
	PartitionStrategy       string // partitionning strategy
}

type ConfigOptions func(options *Options)

// GroupID
func GroupID(id string) ConfigOptions {
	return func(options *Options) {
		options.GroupID = id
	}
}

// Hosts host name
func Hosts(hosts []string) ConfigOptions {
	return func(options *Options) {
		options.Hosts = hosts
	}
}

// Topics set Topic
func Topics(topics ...string) ConfigOptions {
	return func(options *Options) {
		options.Topics = topics
	}
}

// Handler
func Handler(handler func(ctx context.Context, message *sarama.ConsumerMessage) error) ConfigOptions {
	return func(options *Options) {
		options.Handle = handler
	}
}

// Retry for consumers
func Retry(retry bool) ConfigOptions {
	return func(options *Options) {
		options.Retry = retry
	}
}

// SuccessCallback for asyn msg send
func SuccessCallback(callback func(message *sarama.ProducerMessage)) ConfigOptions {
	return func(options *Options) {
		options.ProducerSuccessCallback = callback
	}
}

// ErrorCallback for asyn msg send
func ErrorCallback(callback func(message *sarama.ProducerError)) ConfigOptions {
	return func(options *Options) {
		options.ProducerErrorCallback = callback
	}
}

// PartitionStrategy
func PartitionStrategy(strategy string) ConfigOptions {
	return func(options *Options) {
		options.PartitionStrategy = strategy
	}
}
