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
	"github.com/sg83/go-microservice/article-api/database"
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
	sugar := logger.Sugar()
	sugar.Info("Hello from zap logger")
	defer logger.Sync()

	v := database.NewValidation()

	//Connect to database
	db := database.NewDB(logger)
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
	//re := regexp.MustCompile(`^/tags/(?:[^/]+)/((?:2016|2022)(?:0[1-9]|1[0-2])(?:0[1-9]|[12][0-9]|3[01]))$`)

	//getR.HandleFunc("/tags/{tag:(?:[^/]+)}/{date:((?:2016|2022)(?:0[1-9]|1[0-2])(?:0[1-9]|[12][0-9]|3[01]))$)}", ah.GetTagSummary)
	getR.HandleFunc("/tags/{tag:[a-z]+}", ah.GetTagSummary)

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
		sugar.Info("Starting server on port ", bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			sugar.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	sugar.Info("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
