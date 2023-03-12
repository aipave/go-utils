package gkafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	graceful "github.com/aipave/go-utils/gexit"
	"github.com/sirupsen/logrus"
)

// InitAsyncProducer initializes the asynchronous interface.
// Important:
// If c.Producer.Return.Errors = true, kafka.ErrorCallback must be configured, otherwise there is a risk of panic.
// If c.Producer.Return.Success = true, kafka.SuccessCallback must be configured, otherwise there is a risk of panic.
func InitAsyncProducer(cfg *sarama.Config, opts ...ConfigOptions) (asyncProducer AsyncProducer) {
	// config
	var option Options
	for _, fn := range opts {
		fn(&option)
	}

	logrus.Infof("begin init async producer options:%+v, config:%+v", option, cfg)

	// init Producer
	producer, err := sarama.NewAsyncProducer(option.Hosts, cfg)
	if err != nil {
		panic(fmt.Sprintf("init kafka async producer err:%v hosts:%v config:%+v", err, option.Hosts, cfg))
	}

	// handle error message
	if cfg.Producer.Return.Errors {
		go func() {
			for errMsg := range producer.Errors() {
				go option.ProducerErrorCallback(errMsg)
			}
		}()
	}

	// handle success message
	if cfg.Producer.Return.Successes {
		go func() {
			for successMsg := range producer.Successes() {
				go option.ProducerSuccessCallback(successMsg)
			}
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())

	asyncProducer.Producer = producer
	asyncProducer.ctx = ctx

	graceful.Release(func() {
		cancel()
		err = producer.Close()
		if err != nil {
			logrus.Errorf("close producer:%+v err:%v", producer, err)
		}

		logrus.Infof("kafka producer:%v resource released", option.Hosts)
	})

	logrus.Infof("init async producer success, options:%+v, config:%+v", option, cfg)
	return
}

// PublishAsyncMessage
func PublishAsyncMessage(producer AsyncProducer, topic string, message Message) (err error) {
	if len(topic) == 0 {
		logrus.Errorf("empty topic:%v is unsupported, message:%+v", topic, message)
		return
	}

	if !message.Check() {
		logrus.Errorf("missing producerSeqID or producerTimestamp, topic:%v, message:%+v", topic, message)
		return
	}

	value, _ := json.Marshal(message)

	var producerMsg = &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	select {
	case <-producer.ctx.Done():
		logrus.Warningf("context done, drop topic:%v msg:%+v", topic, message)
		return
	default:
	}

	producer.Producer.Input() <- producerMsg
	logrus.Infof("produced async msg, topic:%v, message:%+v", topic, message)
	return
}

// PublishAsyncRawMessage
func PublishAsyncRawMessage(producer AsyncProducer, producerMessage *sarama.ProducerMessage) (err error) {
	select {
	case <-producer.ctx.Done():
		logrus.Warningf("context done, drop message:%+v", producerMessage)
		return
	default:
	}

	producer.Producer.Input() <- producerMessage
	logrus.Infof("produced async msg, message:%+v", producerMessage)
	return
}
