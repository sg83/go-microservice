package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/data"
	"github.com/sg83/go-microservice/article-api/handlers"
	"go.uber.org/zap"
)

var bindAddress = ":8080"

func main() {

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	//Initialize data validator
	v := data.NewValidation()

	//Connect to database
	db := data.NewDB(logger)
	defer db.Close()

	//Create handlers
	ah := handlers.NewArticles(logger, db, v)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	//Create a new serve mux
	sm := mux.NewRouter()

	//Register handlers for the API's
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/articles/{id:[0-9]+}", ah.Get)
	getR.HandleFunc("/tags/{tag}/{date}", ah.GetTagSummary)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/articles", ah.Create)
	postR.Use(ah.MiddlewareValidateArticle)

	//Create a new server
	s := http.Server{
		Addr:    bindAddress, // configure the bind address
		Handler: ch(sm),      // set the default handler
		ErrorLog: zap.NewStdLog(logger.With(
			zap.String("source", "http-server"),
			zap.String("type", "error-log"),
		)), // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Info("Starting server on port ", zap.String("address", bindAddress))

		err := s.ListenAndServe()
		if err != nil {
			logger.Error("Error starting server", zap.Error(err))
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	logger.Info("Got signal:", zap.Any("signal", sig))

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
