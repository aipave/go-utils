package gkafka

import (
	"context"

	"github.com/Shopify/sarama"
)

// Message
type Message interface {
	Check() bool
}

const (
	ContentTypeJSON      = "json"
	ContentTypeProto     = "protobuf"
	ContentTypeAppCustom = "app/custom"
)

// Header
type Header struct {
	// the ID and CorrelationID fields both provide unique identifiers for messages
	// the ID is used more generally to identify messages
	// the CorrelationID is specifically designed to help correlate related messages within a single request-response flow
	SeqID string `json:"id"` // message sequence
	///> import "github.com/google/uuid" with ID uuid.New() to generate a UUID like "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	///> the TraceID as follows: "f81d4fae-7dec-11d0-a765-00a0c91e6bf6.auth-service.1" mean "`uuid`.`service name`.`the instance identifier`"
	TraceID       string `json:"traceID"`       // used to track a request as it flows through multiple services, and it may change as the request is processed by different services. Each service that handles the request would typically add its own identifier to the TraceID field, creating a trace that shows the path of the request through the system.
	CorrelationID string `json:"correlationID"` // In distributed systems, used to correlate related messages within a single request-response flow, remaining the same throughout the entire request-response flow, regardless of how many services are involved in processing the request.

	Topic       string `json:"topic"`       // for routing messages to the appropriate topic based on their content or destination
	ContentType string `json:"contentType"` // indicating the content type or format of the message payload (e.g. default JSON, Avro, Protobuf, etc.)

	Key KeyInfo `json:"keyInfo"` // This can be useful for ensuring that related messages are sent to the same partition, or for grouping messages based on a shared attribute

	Version string `json:"version"` // Backward compatibility, Graceful migration, Auditability
	Source  string `json:"source"`  // indicating the source of the message, such as the name of the service or application that produced the message, for debugging and tracing msg flows

	Timestamp int64 `json:"timestamp"` // message creation timestamp
}

const (
	KeyTypeUserID     = "key/uid"      // all messages related to that the specific user are sent to the same partition, so they can be processed in order.
	KeyDeviceID       = "key/device"   // In an e-commerce platform, use the order ID as the Key for messages related to a specific order
	KeyTypeOrderID    = "key/order"    // use the product ID as the Key for messages related to a specific product
	KeyTypeProductID  = "key/product"  // In a geolocation-based system, ensures that all messages related to that location are sent to the same partition
	KeyTypeLocationID = "key/location" // In an IoT platform, ensures that all messages related to that device are sent to the same partition
)

type KeyInfo struct {
	KeyType    string `json:"keyType"`
	KeyContent string `json:"keyContent"`
}

// Check 消息头校验
func (h Header) Check() bool {
	if len(h.SeqID) == 0 || h.Timestamp == 0 {
		return false
	}

	return true
}

// AsyncProducer 异步生产者
type AsyncProducer struct {
	ctx      context.Context
	Producer sarama.AsyncProducer
}
