package gmongo

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/alyu01/go-utils/gracexit"
    "github.com/sirupsen/logrus"
)

// InitMongo initial mongo connect
// gmongo.MongoConfig or URI String
func InitMongo(cc any) (client *mongo.Client) {
    var uri string
    switch cc := cc.(type) {
    case MongoConfig:
        uri = cc.URI()
    case string:
        uri = cc
    default:
        panic(fmt.Errorf("unkown type:%v", cc))
    }

    var err error
    logrus.Infof("connecting to addr:%v", uri)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        panic(fmt.Errorf("connecting to mongo addr:%v err:%v", uri, err))
    }

    gracexit.Release(func() {
        if err = client.Disconnect(context.Background()); err != nil {
            logrus.Errorf("disconnect from mongo addr:%v err:%v", uri, err)
        }

        logrus.Infof("mongo addr:%v resource released.", uri)
    })

    logrus.Infof("connected to mongo addr:%v done", uri)
    return
}
