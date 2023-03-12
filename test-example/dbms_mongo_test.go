package test_example

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	algorithms "github.com/aipave/go-utils/algorithms"
	"github.com/aipave/go-utils/dbms/gmongo"
	"github.com/aipave/go-utils/gcast"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"

	//"github.com/aipave/go-utils/dbms/gredis"
	"github.com/aipave/go-utils/ginfos"
	"gopkg.in/yaml.v3"
)

var configFile = flag.String("f", "config.yaml", "the config file for dbms test")

func GetConfig() Config {
	return globalConfig
}

var globalConfig Config

type Config struct {
	RWMysql string             `yaml:"RWMysql"`
	Mongo   gmongo.MongoConfig `yaml:"Mongo"`
	//Redis   gredis.RedisConfig `yaml:"Redis"`
}

func init() {
	ginfos.Version()

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		//panic("open file err:" + err.Error())
		logrus.Errorf("open file err %v", err)
	}
	err = yaml.Unmarshal(data, &globalConfig)
	if err != nil {
		//panic("open file err:" + err.Error())
		logrus.Errorf("open file err %v", err)
	}

}

const (
	MongoDbTest         = "test"
	UserTbInMongoDbTest = "users"
)

var mongoCli *mongo.Client

func TestMongoDb(t *testing.T) {
	t.Log(*configFile)
	t.Log(GetConfig().Mongo)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		mongoCli.Disconnect(context.Background())
	}()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?retryWrites=false"))
	if err != nil {
		//
		t.Fatalf("err: %v", err)

	}

	mongoCli = client // 将 client 赋值给 mongoCli
	t.Logf("connect ok|%v", mongoCli)

	//t.Log(DeleteAll(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{})) //ok
	for cnt := 10; cnt >= 0; cnt-- {
		t.Log((Insert(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{
			Seq:        algorithms.Snowflake.NextOptStreamID(),
			Uid:        100000 + rand.Int63()%900000,
			Age:        22 + rand.Int31()%45,
			Gender:     gcast.ToBool(rand.Int() % 2),
			UpdateTime: time.Now().UnixMilli() / 1000,
		})))

		// ok
		//t.Log(Delete(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{
		//    Uid: 293410,
		//}))

		// ok
		//t.Log(Search(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{
		//    Uid: 9410,
		//}))

	}
	InsertTransaction(context.Background()) // not ok
	t.Log(FindMore(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{}))

}

type User struct {
	//Key primitive.ObjectID `bson:"_id, omitempty"` /* If no id value is provided, mongo db will automatically generate an object id and assign it to the id field.*/
	Seq        int64 `bson:"id"`
	Uid        int64 `bson:"uid"`
	Age        int32 `bson:"age"`
	Gender     bool  `bson:"gender"`
	UpdateTime int64 `bson:"timestamp"`
}

func Insert(ctx context.Context, dbi, table string, user User) bool {
	///> get database and collection
	db := mongoCli.Database(dbi)
	_, err := db.Collection(table).InsertOne(ctx, user)
	if err != nil {
		logrus.Errorf("fail, %v", err)
		return false
	}
	logrus.Infof("%v ok|%v", ginfos.FuncName(), user)
	return true

}

func Search(ctx context.Context, dbi, table string, useri User) (user User, ok bool) {
	err := mongoCli.Database(dbi).Collection(table).FindOne(ctx,
		bson.M{
			"update_time": bson.M{"$gt": time.Now().UnixMilli() / 1000},
		}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return
	}
	if err != nil {
		logrus.Errorf("%v|err", ginfos.FuncName())
	}
	ok = true
	return
}

func FindMore(ctx context.Context, dbi, table string, useri User) (users []User) {
	cursor, err := mongoCli.Database(dbi).Collection(table).Find(ctx,
		bson.M{
			"timestamp": bson.M{"$gt": time.Now().UnixMilli()/1000 - 100},
			"gender":    true,
		},
		// Error : use the SetLimit method first and then the SetSkip method
		// The number of documents counted from the beginning in the query results are sorted by insertion time in default
		// e. set skip 200 means skip the first 200 documents
		options.Find().SetSkip(30),

		options.Find().SetLimit(3),

		options.Find().SetSort(bson.D{{"age", -1}}), // 1 means ascending order, -1 means descending order
	)

	// Since the result set may contain a large amount of data,
	// it needs to be closed explicitly to avoid resource leaks,
	// otherwise it will always occupy memory.
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			break
		}
		users = append(users, user)

	}

	return
}

func Delete(ctx context.Context, dbi, table string, useri User) (ok bool) {
	_, err := mongoCli.Database(dbi).Collection(table).DeleteOne(ctx,
		bson.M{
			"uid": useri.Uid,
		})
	if err != nil {
		logrus.Errorf("%v|err", ginfos.FuncName())
	}
	ok = true
	return
}

func DeleteAll(ctx context.Context, dbi, table string, useri User) (ok bool) {
	_, err := mongoCli.Database(dbi).Collection(table).DeleteMany(ctx,
		bson.M{})
	if err != nil {
		logrus.Errorf("%v|err", ginfos.FuncName())
	}
	ok = true
	return
}

func InsertTransaction(ctx context.Context) (err error) {
	sessionErr := mongoCli.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// start transaction
		_, err = sessionContext.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			Insert(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{
				Seq:        algorithms.Snowflake.NextOptStreamID(),
				Uid:        666888,
				Age:        22 + rand.Int31()%45,
				Gender:     gcast.ToBool(rand.Int() % 2),
				UpdateTime: time.Now().UnixMilli() / 1000,
			})

			return nil, nil

		})
		if err != nil {
			logrus.Fatalf("%v--|%v", ginfos.FuncName(), err)
			return err
		}
		return nil
	})
	if sessionErr != nil {
		logrus.Fatalf("%v--|%v", ginfos.FuncName(), sessionErr)
		return sessionErr
	}
	return nil

}

func InsertTransaction2(ctx context.Context) {
	session, err := mongoCli.StartSession()
	if err != nil {
		logrus.Fatalf("%v-|%v", ginfos.FuncName(), err)
	}
	defer session.EndSession(ctx)

	cb := func(sessCtx mongo.SessionContext) (interface{}, error) {
		Insert(context.Background(), MongoDbTest, UserTbInMongoDbTest, User{
			Seq:        algorithms.Snowflake.NextOptStreamID(),
			Uid:        666888,
			Age:        22 + rand.Int31()%45,
			Gender:     gcast.ToBool(rand.Int() % 2),
			UpdateTime: time.Now().UnixMilli() / 1000,
		})

		return nil, nil
	}

	result, err2 := session.WithTransaction(context.Background(), cb)
	if err != nil {
		logrus.Fatalf("%v--|%v", ginfos.FuncName(), err2)
	}
	logrus.Info(result)

}
