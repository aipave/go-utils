package gkafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	grace "github.com/aipave/go-utils/gexit"
	ginfos "github.com/aipave/go-utils/ginfos"
	"github.com/sirupsen/logrus"
)

// InitSyncProducer
func InitSyncProducer(opts ...ConfigOptions) (producer sarama.SyncProducer) {
	var option Options
	for _, fn := range opts {
		fn(&option)
	}

	logrus.Infof("begin init sync producer option:%+v", option)

	// init Producer
	var err error
	var c = sarama.NewConfig()
	c.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer(option.Hosts, c)
	if err != nil {
		panic(fmt.Sprintf("init kafka sync producer err:%v, hosts:%v, config:%+v", err, option.Hosts, c))
	}

	grace.Release(func() {
		err = producer.Close()
		if err != nil {
			logrus.Errorf("close producer:%+v err:%v", producer, err)
		}

		logrus.Infof("producer:%v resource released", option.Hosts)
	})

	logrus.Infof("init sync producer success, option:%+v, config:%v", option, c)
	return producer
}

// PublishSyncMessage
func PublishSyncMessage(producer sarama.SyncProducer, topic string, message Message) (err error) {
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

	var partition, offset = int32(0), int64(0)
	partition, offset, err = producer.SendMessage(producerMsg)
	if err != nil {
		logrus.Errorf("send message failed, topic:%v partition:%v offset:%v err:%v value:%v", topic, partition, offset, err, string(value))
	} else {
		logrus.Infof("send message success, topic:%v partition:%v offset:%v value:%v", topic, partition, offset, string(value))
	}

	return
}

// SyncPublish
func SyncPublish(producer sarama.SyncProducer, topic string, message []byte) (err error) {
	if len(topic) == 0 {
		logrus.Errorf("empty topic:%v is unsupported, message:%+v", topic, message)
		return
	}

	var producerMsg = &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	var partition, offset = int32(0), int64(0)
	partition, offset, err = producer.SendMessage(producerMsg)
	if err != nil {
		logrus.Errorf("send message failed, topic:%v partition:%v offset:%v err:%v value:%v", topic, partition, offset, err, string(message))
	} else {
		logrus.Infof("send message success, topic:%v partition:%v offset:%v value:%v", topic, partition, offset, string(message))
	}

	return
}

// PublishDelayMessage Delayed Queue
// unique -- Unique identifier for the message, ensuring uniqueness within the Topic dimension is sufficient
// deadline -- Deadline for the delay to end, in seconds
// topic -- Delayed queue TOPIC (Message feedback queue)
// msg -- Delayed message content
func PublishDelayMessage(producer sarama.SyncProducer, unique string, deadline int64, topic string, msg []byte) (err error) {
	message := ginfos.JsonStr(map[string]interface{}{
		"unique":   unique,
		"deadline": deadline,
		"topic":    topic,
		"msg":      string(msg),
	})

	var producerMsg = &sarama.ProducerMessage{
		Topic: "delay_queue_svr",
		Value: sarama.StringEncoder(message),
	}

	var partition, offset = int32(0), int64(0)
	partition, offset, err = producer.SendMessage(producerMsg)
	if err != nil {
		logrus.Errorf("send message failed, topic:%v partition:%v offset:%v err:%v value:%v", topic, partition, offset, err, message)
	} else {
		logrus.Infof("send message success, topic:%v partition:%v offset:%v value:%v", topic, partition, offset, message)
	}

	return
}
