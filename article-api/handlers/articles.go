package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/database"
	"go.uber.org/zap"
)

type Articles struct {
	l  *zap.Logger
	db *database.ArticlesDb
}

func NewArticles(l *zap.Logger, db *database.ArticlesDb) *Articles {
	return &Articles{l, db}
}

func (a *Articles) Get(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	a.l.Info("Get article", zap.String("id", vars["id"]))

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		a.l.Fatal("Could not convert id to int")
		return
	}
	w.Header().Add("Content-Type", "application/json")

	article, err := a.db.GetArticleByID(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		w.WriteHeader(http.StatusInternalServerError)
		database.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = database.ToJSON(article, w)
	if err != nil {
		// we should never be here but log the error just incase
		a.l.Error("Unable to serialize product", zap.String(" Error: ", err.Error()))
		return
	}
}

func (a *Articles) Create(w http.ResponseWriter, r *http.Request) {

	a.l.Info("Create article")
}

func (a *Articles) GetTagSummary(w http.ResponseWriter, r *http.Request) {
	a.l.Info("Get tag summary")
}
