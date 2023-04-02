package main

import (
	"context"
	"email-project/base"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	ruslogger := logrus.New()
	ruslogger.SetLevel(logrus.ErrorLevel)

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		ruslogger.Error("error loading .env file")
	}

	// Access environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbClusterInfo := os.Getenv("DB_CLUSTER_INFO")

	// Use environment variables
	dbUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s", dbUsername, dbPassword, dbClusterInfo)
	// ...

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// db setup
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		ruslogger.Error("unable to setup new mongodb client")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		ruslogger.Error("unable to connect to mongodb client")
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		ruslogger.Error("unable to ping mongodb client")
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		ruslogger.Error("unable to list mongodb client")
	}
	usersCollection := client.Database("authenticate-db").Collection("email-project")

	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	// insert the bson object using InsertOne()
	result, err := usersCollection.InsertOne(context.TODO(), user)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}
	// display the id of the newly inserted object
	fmt.Println(result.InsertedID)

	fmt.Println(databases)

	var s base.Service
	{
		s = base.NewService(*client)
		s = base.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = base.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
