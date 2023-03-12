package gkafka

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Shopify/sarama"
	grace "github.com/aipave/go-utils/gexit"
	"github.com/sirupsen/logrus"
)

const (
	PartitionStrategyRange      = "range"
	PartitionStrategyRoundRobin = "roundrobin"
	PartitionStrategySticky     = "sticky"
)

type consumer struct {
	session sarama.ConsumerGroupSession
	message *sarama.ConsumerMessage
}

var topicChannel = make(map[string]chan *consumer)

// InitKafkaConsumerGroup(Host(xxx), Topics(xxx), Handler(xxx))
func InitKafkaConsumerGroup(opts ...ConfigOptions) {
	var option Options
	for _, fn := range opts {
		fn(&option)
	}

	logrus.Infof("start init consumer group addr:%v topics:%v", option.Hosts, option.Topics)

	// kafka config
	cfg := sarama.NewConfig()
	// defaul sticky strategy
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	cfg.Consumer.Group.Session.Timeout = 10 * time.Second
	cfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	cfg.Consumer.IsolationLevel = sarama.ReadCommitted
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Version = sarama.V2_5_0_0

	// set partition
	if option.PartitionStrategy == PartitionStrategyRange {
		cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	} else if option.PartitionStrategy == PartitionStrategyRoundRobin {
		cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	}

	// set consumer group
	var groupID = option.GroupID
	if len(groupID) == 0 {
		groupID = defaultPath() // default is process_name
	}
	consumerGroup, err := sarama.NewConsumerGroup(option.Hosts, groupID, cfg)
	if err != nil {
		panic(fmt.Sprintf("init kafka consumer err:%v host:%v topics:%v", err, option.Hosts, option.Topics))
	}

	if cfg.Consumer.Return.Errors {
		go func() {
			for errMsg := range consumerGroup.Errors() {
				if option.ConsumerErrorCallback != nil {
					option.ConsumerErrorCallback(errMsg)
				}
			}
		}()
	}

	// init queue
	for _, topic := range option.Topics {
		var topicChan = make(chan *consumer, 100)
		go topicConsumer(option, topicChan)
		topicChannel[topic] = topicChan
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			err = consumerGroup.Consume(ctx, option.Topics, &consumerGroupHandler{opts: option, handle: option.Handle})
			if err != nil {
				logrus.Errorf("consumer returned with err:%v addr:%v topics:%v", err, option.Hosts, option.Topics)
			}

			if ctx.Err() != nil {
				logrus.Warningf("ctx done: %v", ctx.Err())
				return
			}
		}
	}()

	grace.Close(func() {
		cancel()
		err = consumerGroup.Close()
		if err != nil {
			logrus.Errorf("consumer group close err:%v host:%v topics:%v", err, option.Hosts, option.Topics)
		}

		logrus.Infof("kafka addr:%v consumer closed", option.Hosts)
	})

	logrus.Infof("init consumer group success addr:%v topics:%v", option.Hosts, option.Topics)
}

func defaultPath() string {
	pwd, _ := os.Executable()
	_, exec := filepath.Split(pwd)
	return exec
}

type consumerGroupHandler struct {
	opts   Options
	handle func(ctx context.Context, message *sarama.ConsumerMessage) error
}

// Setup hooks before consumeClaim
func (c consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup hooks after consumeClaim
func (c consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim
func (c consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg == nil {
			session.MarkMessage(msg, "")
			logrus.Warningf("handle nil msg, continue")
			continue
		}

		if topicChan, ok := topicChannel[msg.Topic]; ok {
			topicChan <- &consumer{session: session, message: msg}
		} else {
			err := c.opts.Handle(session.Context(), msg)
			if err != nil {
				logrus.Warningf("default|msg=%v|err=%v", string(msg.Value), err)
			}
		}
	}

	return nil
}

func topicConsumer(opts Options, consumerChan chan *consumer) {
	for consumerMessage := range consumerChan {
		session, msg := consumerMessage.session, consumerMessage.message

		// handle msg
		err := opts.Handle(session.Context(), msg)
		if err != nil {
			logrus.Errorf("handle msg:%+v, timestamp:%v, header:%+v, topic:%v, partition:%v, offset:%v, err:%v", string(msg.Value), msg.Timestamp, msg.Headers, msg.Topic, msg.Partition, msg.Offset, err)
		}

		// consumer success or drop failed message
		if err == nil || !opts.Retry {
			session.MarkMessage(msg, "")
		} else {
			// TODO: handle retry
		}
	}

}
