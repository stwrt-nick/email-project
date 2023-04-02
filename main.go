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
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ruslogger := logrus.New()
	ruslogger.SetLevel(logrus.ErrorLevel)

	// db setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		ruslogger.Error("unable to connect to mongodb")
	}

	// Check the connection
	err = db.Ping(ctx, nil)
	if err != nil {
		ruslogger.Error("unable to connect to mongodb")
	}

	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			ruslogger.Error("db disconnected")
		}
	}()

	var s base.Service
	{
		s = base.NewService(*db)
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
