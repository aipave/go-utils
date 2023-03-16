package test_example

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/aipave/go-utils/algorithms"
	"github.com/aipave/go-utils/gcast"
	"github.com/aipave/go-utils/gexit"
	"github.com/aipave/go-utils/ginfos"
	"github.com/aipave/go-utils/gkafka"
	"github.com/aipave/go-utils/glogs/glogrus"
	"github.com/hashicorp/go-uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// /> "github.com/Shopify/sarama"
func TestKafkaRead(t *testing.T) {
	// Set up a Kafka reader
	brokers := []string{"localhost:9092", "localhost:9093"} // replace with your external address
	topic := "signsvr"
	partition := 0

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: partition,
	})

	// Read messages from Kafka
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message from Kafka:", err.Error())
			break
		}

		fmt.Printf("Received message: %s\n", string(msg.Value))
	}
}

// /> "github.com/Shopify/sarama"
func TestKafkaConnect01(t *testing.T) {
	//brokers := []string{"localhost:9093"} ///< for docker_compose_kafka-unuseful.yaml
	brokers := []string{"localhost:9092", "localhost:9093"} ///< for docker_compose_kafka_cluster.yaml

	config := sarama.NewConfig()
	config.ClientID = "TestKafkaConfig"
	config.Consumer.Return.Errors = true

	//create a new consumer
	var err error
	var consumer sarama.Consumer
	if consumer, err = sarama.NewConsumer(brokers, config); err != nil {
		logrus.Fatalf("Error creating consumer:%v", err)
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			logrus.Errorf("Error closing consumer:%v", err)
		}
	}()

	partition := int32(0)
	offset := sarama.OffsetNewest

	// Start consuming from the specified partition
	var partitionConsumer sarama.PartitionConsumer
	partitionConsumer, err = consumer.ConsumePartition("signsvr", partition, offset)
	if err != nil {
		logrus.Fatalf("Error creating partition consumer:%v", err)
	}

	defer func() {
		if err = partitionConsumer.Close(); err != nil {
			logrus.Errorf("Error closing partitionConsumer:%v", err)
		}
	}()

	// Set up a signal channel to handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consume messages and handle them in a separate goroutine
	go func() {
		for msg := range partitionConsumer.Messages() {
			logrus.Infof("Received message with key %s and value %s\n", string(msg.Key), string(msg.Value))
		}
	}()

	// Wait for a signal to gracefully shutdown the consumer
	<-signals
}

var kafkaConfigFile = flag.String("f", "config.yaml", "the config file for dbms test")

func GetKafkaConfig() KafkaConfig {
	return kafkaConfig
}

var kafkaConfig KafkaConfig

type KafkaConfig struct {
	Log struct {
		Level string `yaml:"Level"`
		Path  string `yaml:"Path"`
		Mode  string `yaml:"Mode"`
	} `yaml:"Log"`

	Kafka string `yaml:"Kafka"`
}

func init() {
	data, err := ioutil.ReadFile(*kafkaConfigFile)
	if err != nil {
		logrus.Errorf("open file err %v", err)
	}
	err = yaml.Unmarshal(data, &kafkaConfig)
	if err != nil {
		logrus.Errorf("unmarshal file err|%v", err)
	}

	glogrus.Init(glogrus.WithAlerUrl("https://open.feishu.cn/open-apis/bot/v2/hook/2f1dc72c-8d2d-4641-bd95-31bbd6fcd2c7"))
	glogrus.MustSetLevel(GetKafkaConfig().Log.Level)
}

const (
	TopicCoinOption = "coinOption" // charge, add, decrease
	TopicFlowFlow   = "followFlow" // follow, unfollow
)

func TestKafkaConsumerGroup(t *testing.T) {
	// Producer
	go func() {
		var syncProducer sarama.SyncProducer
		syncProducer = gkafka.InitSyncProducer(gkafka.Hosts(strings.Split(GetKafkaConfig().Kafka, ",")))
		logrus.Infof("init kafka host:%v", GetKafkaConfig().Kafka)

		tick := time.Tick(1 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())
		gexit.Close(cancel)
	Loop:
		for {
			select {
			case <-tick:
				uu, _ := uuid.GenerateUUID()
				sid := algorithms.Snowflake.NextOptStreamID()
				var msg *UserCoinOption = &UserCoinOption{
					Header: gkafka.Header{
						SeqID:         gcast.ToString(sid),
						TraceID:       fmt.Sprintf("%v.%v", uu, ginfos.FuncName()),
						CorrelationID: uu,
						Topic:         TopicCoinOption,
						ContentType:   gkafka.ContentTypeJSON,
						Key: gkafka.KeyInfo{
							KeyType:    gkafka.KeyTypeProductID,
							KeyContent: "my_app",
						},
						Source:    ginfos.Runtime.IP(),
						Timestamp: time.Now().Unix(),
					},
					SID:     sid,
					UID:     123456,
					Type:    100,
					Amount:  int64(1) + rand.Int63()%99,
					Balance: 100,
				}
				if err := gkafka.PublishSyncMessage(syncProducer, TopicCoinOption, msg); err != nil {
					logrus.Errorf("publish error|msg=%v, err:%v", msg, err)
				} else {
					logrus.Infof("[producer]push ok|msg=%v", msg)
				}
			case <-ctx.Done():
				logrus.Errorf("publish error")
				break Loop

			}
		}

	}()

	// Consumer
	go func() {
		gkafka.InitKafkaConsumerGroup(gkafka.Hosts(strings.Split(GetKafkaConfig().Kafka, ",")),
			gkafka.Retry(true),
			gkafka.Topics(
				TopicCoinOption,
				TopicFlowFlow,
			),
			gkafka.Handler(kafkaHander))
	}()

	gexit.Wait()
}

func kafkaHander(ctx context.Context, consumerMsg *sarama.ConsumerMessage) (err error) {
	switch consumerMsg.Topic {
	case TopicCoinOption:
		err = handleCoinOption(ctx, consumerMsg.Value)
	case TopicFlowFlow:
	default:
		logrus.Warningf("unknown topic:%v,msg:%v", consumerMsg.Topic, consumerMsg.Value)

	}
	return nil
}

type UserCoinOption struct {
	gkafka.Header
	SID     int64 `json:"sid"`
	UID     int64 `json:"uid"`
	Type    int64 `json:"type"`
	Amount  int64 `json:"amount"`
	Balance int64 `json:"balance"`
}

func (msg *UserCoinOption) RandMarshal() (err error) {
	msg = &UserCoinOption{}
	uu, _ := uuid.GenerateUUID()
	msg.Header = gkafka.Header{
		SeqID:         gcast.ToString(algorithms.Snowflake.NextOptStreamID()),
		TraceID:       fmt.Sprintf("%v.%v", uu, ginfos.FuncName()),
		CorrelationID: uu,
		Topic:         TopicCoinOption,
		ContentType:   gkafka.ContentTypeJSON,
		Key:           gkafka.KeyInfo{},
		Source:        ginfos.Runtime.IP(),
		Timestamp:     time.Now().Unix(),
	}
	return nil
}

func (msg *UserCoinOption) Unmarshal(b []byte) (err error) {
	err = json.Unmarshal(b, msg)
	return err
}

func handleCoinOption(ctx context.Context, b []byte) (err error) {
	var msg *UserCoinOption = &UserCoinOption{}
	err = msg.Unmarshal(b)
	if err != nil {
		logrus.Errorf("unmarshal err:%v", b)
		return err
	}

	logrus.Infof("[consumer]marshal ok:%v", msg)
	return nil
}
