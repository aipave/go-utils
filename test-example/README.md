# test example

## 1. dbms mongo test

### cursor.Close(ctx) ?

### ....

## 2. dbms redis test

### 2.1 Redis ACID

Strictly speaking, redis has no transactions. Because transactions must have four characteristics: atomicity,
consistency, isolation, and durability. Then redis cannot do these four points, but only has some of these
characteristics. The redis transaction is a pseudo transaction and does not support rollback.

Redis transaction mechanism can ensure consistency and isolation, but cannot guarantee persistence. It has some
atomicity, but does not support rollback.

(1) Atomicity:(multi~exec:ok? exec-ing:ok, after exec:not ok)

After multi, Such as insufficient memory, wrong command name, redis will report an error and record the error. At this
time, the redisClient can continue to submit command operations; when executing exec, redis will refuse to execute all
submitted command operations, and return the result of transaction failure nil

If a failure occurs when EXEC is executed, If Redis has enabled AOF logging, only part of the transactional operations
will be recorded in the AOF log. You need to use the redis-check-aof tool to check the AOF log file, which can remove
unfinished transactional operations from the AOF file. The atomicity of the transaction is ensured.

After exec, The data types of commands and operations do not match, but the Redis instance does not detect an error.
After executing the EXEC command, Redis actually executes these instructions, and an error occurs. At this time, the
transaction will not be rolled back, but the commands in the transaction queue will still be executed. The atomicity of
the transaction cannot be guaranteed.

(2) Consistency(multi~exec:ok, exec-ing:ok, after exec:ok)

Before EXEC command execution: If there is an error during enqueueing, the transaction will be abandoned, ensuring
consistency.

If a failure occurs during EXEC command execution: In RDB mode, the RDB snapshot is not executed during transaction
execution, and the transaction result is not saved in RDB; in AOF mode, you can use the redis-check-aof tool to remove
unfinished transactional operations from the AOF file. Consistency can be ensured.

After EXEC command execution: If there is an error during actual execution, the incorrect command will not be executed,
and the correct commands will be executed normally, ensuring consistency.

(3) Isolation(multi~exec:ok, after exec:ok):

Before the execution of the EXEC command, isolation needs to be guaranteed through the WATCH mechanism. This is because
before the execution of the EXEC command, other client commands can be executed, and relevant variables may be modified;
but the WATCH mechanism can be used to monitor relevant variables. Once the relevant variables are modified, the
transaction fails and returns after the EXEC command; ensuring isolation.

After the execution of the EXEC command, isolation can be guaranteed. This is because Redis is single-threaded, and the
commands in the transaction queue and the commands from other clients can only be executed in sequence. Therefore, it
ensures isolation.

(4) Durability(multi~exec:0, after exec:0)

If Redis does not use RDB or AOF, transaction persistence does not exist.

In RDB mode, if an instance crashes after a transaction is executed and before the next RDB snapshot is taken, data loss
may occur. In this case, the data modified by the transaction cannot be guaranteed to be persistent.

In AOF mode, because all three configuration options (no, everysec, and always) of AOF mode have the possibility of data
loss, the transaction's persistence attribute cannot be guaranteed either.

### 2.2 Advantage of lua script

Reducing network overhead: multiple requests can be sent at once in the form of a script, reducing network latency.

Atomic operation: Redis executes the entire script as a whole and no other requests will be inserted in the middle.
There is no need to worry about race conditions during script execution.

Reusable: The script sent by the client will permanently exist in Redis, so other clients can reuse the same script
without needing to write the same logic in their code.

### 2.3 example for redis transaction

```shell
redis> WATCH "name" ### Used to monitor the redis variable value, before the command exec; the data in redis may be modified by other client commands
OK
redis> MULTI  ### The client executing this command switches from a non-transactional state to a transactional state 
OK          ### After multi starts a transaction, it is not a special command such as watch, exec, discard, multi; the command of the client will not be executed immediately, but put into a transaction queue
redis> SET "name" "lwl" ### if name has been changed by others
QUEUED
redis> EXEC ### If an exec command is received, the commands in the transaction queue will be executed. If it is discard, the transaction is discarded
(nil)  ### will not return ok
```

If there is an error in the command enqueuing process, such as using a non-existing command, the transaction queue will
be rejected for execution

An exception occurs during the execution of the transaction, such as the data type of the command and the operation do
not match, and the command in the transaction queue will continue Execution continues until all commands are executed.
will not roll back

### 2.4 Example grab envelope

The set object in Redis is unordered and unique. Set collection is implemented by an integer collection or a dictionary,
and the complexity of adding, deleting, and searching is basically considered as O(1). The maximum number of objects
that can be stored is 2^32 - 1 (4294967295).

To keep track of the users who have participated in an activity, we can use a set collection. Before a user participates
in the activity, we can check if the user is in the set collection. If not, the user can grab the red envelope.

If users can participate multiple times, we can use a hash object. The key stores the user object, and the value stores
the number of times the user participated. We can use the INCR atomic operation to increase the value, and if the
returned value is greater than the upper limit, it means the user has reached the maximum participation limit.

Two instructions in the set collection are very suitable for scenarios such as grabbing red envelopes and lottery:

SPOP key [count]: removes and returns a random element from the set
SRANDMEMBER key [count]: returns one or more random numbers from the set; need to call SREM again to remove them
All the red envelopes can be added to the set using SADD, and then the corresponding red envelope can be obtained
through the random command.

If there are empty options such as "thank you for your patronage", corresponding invalid red envelopes or prizes can be
generated and added to the set or list.

Grabbing red envelopes usually has a time limit, which can be perfectly solved by using the expiration time of the Redis
key.

### 2.5 tips

All Redis data types including lists, sets, hashes, and sorted sets (zsets) can have a TTL (time-to-live) set on them.
This means that after the specified time has elapsed, the data will be automatically deleted by Redis. The TTL can be
set using the EXPIRE command, which takes a key and a time in seconds as arguments. The EXPIREAT command can also be
used to set the expiry time as a Unix timestamp.
```
127.0.0.1:6379> expire envelope:1633685538140921856 86400
(integer) 1
127.0.0.1:6379> ttl envelope:1633685538140921856
(integer) 89

```

```redis-cli
script load "return redis.call('set',KEYS[1],ARGV[1])" // Generates a digest for the current Lua script, and returns the digest

evalsha "c686f316aaf1eb01d5a4de1b0b63cd233010e63d" 1 address china // Executes the command "set address china" using the digest of the Lua script

get address // Retrieves "address" to confirm whether the previous script was executed successfully

script exists "c686f316aaf1eb01d5a4de1b0b63cd233010e63d" // Checks whether the Lua script digest exists

script flush // Clears all the Lua script cache

script exists "c686f316aaf1eb01d5a4de1b0b63cd233010e63d" // After the cache is cleared, the digest no longer exists.

```

# 3. kafka cluter test
## 3.1 confluentinc/cp-kafka
### 3.1.1 install docker 
```bash
docker pull confluentinc/cp-kafka
docker pull confluentinc/cp-zookeeper
```
### 3.1.2 docker-compose file
```yaml
---
version: '2'
services:
  cp-zookeeper-1:
    user: root
    restart: always
    container_name: cp-zookeeper-1
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - 22181:2181

  cp-zookeeper-2:
    user: root
    restart: always
    container_name: cp-zookeeper-2
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - 32181:2181
  cp-zookeeper-3:
    user: root
    restart: always
    container_name: cp-zookeeper-3
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - 42181:2181

  kafka-1:
    user: root
    restart: always
    container_name: kafka-1
    image: confluentinc/cp-kafka:latest
    depends_on:
      - cp-zookeeper-1
      - cp-zookeeper-2
      - cp-zookeeper-3

    ports:
      - 29092:29092
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: cp-zookeeper-1:2181,cp-zookeeper-2:2181,cp-zookeeper-3:2181/kafkas-a
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:9092,PLAINTEXT_HOST://localhost:29092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
  kafka-2:
    user: root
    restart: always
    container_name: kafka-2
    image: confluentinc/cp-kafka:latest
    depends_on:
      - cp-zookeeper-1
      - cp-zookeeper-2
      - cp-zookeeper-3
    ports:
      - 39092:39092
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_ZOOKEEPER_CONNECT: cp-zookeeper-1:2181,cp-zookeeper-2:2181,cp-zookeeper-3:2181/kafkas-a
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:9092,PLAINTEXT_HOST://localhost:39092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
  kafka-3:
    user: root
    restart: always
    container_name: kafka-3
    image: confluentinc/cp-kafka:latest
    depends_on:
      - cp-zookeeper-1
      - cp-zookeeper-2
      - cp-zookeeper-3
    ports:
      - 49092:39092
    environment:
      KAFKA_BROKER_ID: 3
      KAFKA_ZOOKEEPER_CONNECT: cp-zookeeper-1:2181,cp-zookeeper-2:2181,cp-zookeeper-3:2181/kafkas-a
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-3:9092,PLAINTEXT_HOST://localhost:39092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

```
### 3.1.3 start docker cluster
```yaml
$ docker-compose  -f docker_compose_cp-kafka-cluster.yaml up -d
[+] Running 6/6
 ⠿ Container cp-zookeeper-2  Started                                                                                                                              0.4s
 ⠿ Container cp-zookeeper-1  Started                                                                                                                              0.4s
 ⠿ Container cp-zookeeper-3  Started                                                                                                                              0.4s
 ⠿ Container kafka-2         Started                                                                                                                              0.5s
 ⠿ Container kafka-1         Started                                                                                                                              0.7s
 ⠿ Container kafka-3         Started    
```

### 3.1.4 go test demo
config in yaml 
```yaml
Kafka: localhost:29092,localhost:39092,localhost:49092 #for cp/kafka
```
code main body
```go
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

```
### 3.1.5 result
```
15039-03-15 10:36:11.291 [/Users/xxx/mod/gomod/pkg/mod/github.com/aipave/go-utils@v0.0.58/gkafka/producer_with_sync.go:67 gkafka.PublishSyncMessage] [INFO] send message success, topic:coinOption partition:0 offset:688 value:{"id":"1635832233163300864","traceID":"3340cb86-464b-790f-92bc-f53a3057b6b1.func1","correlationID":"3340cb86-464b-790f-92bc-f53a3057b6b1","topic":"coinOption","contentType":"json","keyInfo":{"keyType":"key/product","keyContent":"my_app"},"version":"","source":"172.17.17.153","timestamp":1678847771,"sid":1635832233163300864,"uid":123456,"type":100,"amount":33,"balance":100}
15039-03-15 10:36:11.292 [/Users/xxx/mod/go-utils/test-example/gkafka_test.go:184 command-line-arguments.TestKafkaConsumerGroup.func1] [INFO] [producer]push ok|msg=&{{1635832233163300864 3340cb86-464b-790f-92bc-f53a3057b6b1.func1 3340cb86-464b-790f-92bc-f53a3057b6b1 coinOption json {key/product my_app}  172.17.17.153 1678847771} 1635832233163300864 123456 100 33 100}
15039-03-15 10:36:11.292 [/Users/xxx/mod/go-utils/test-example/gkafka_test.go:259 command-line-arguments.handleCoinOption] [INFO] [consumer]marshal ok:&{{1635832233163300864 3340cb86-464b-790f-92bc-f53a3057b6b1.func1 3340cb86-464b-790f-92bc-f53a3057b6b1 coinOption json {key/product my_app}  172.17.17.153 1678847771} 1635832233163300864 123456 100 33 100}

```