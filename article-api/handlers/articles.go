package handlers

import (
	"net/http"
	"strconv"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/database"
	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

// KeyArticle is a key used for the Article object in the context
type KeyArticle struct{}

type Articles struct {
	l  *zap.Logger
	db database.ArticlesData
	v  *database.Validation
}

func NewArticles(l *zap.Logger, db database.ArticlesData, v *database.Validation) *Articles {
	return &Articles{l, db, v}
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

// Create handles POST requests to add new products
func (a *Articles) Create(w http.ResponseWriter, r *http.Request) {

	a.l.Info("Create article", zap.Any("article:", r.Context().Value(KeyArticle{})))

	// fetch the article from the context
	//&article := r.Context().Value(KeyArticle{}).(*models.Article)
	article, ok := r.Context().Value(KeyArticle{}).(*models.Article)
	if !ok {
		// handle the case where the value is not of the expected type
		fmt.Println("fetching object from context")
		return
	}

	a.l.Info("Inserting ", zap.Any("article: ", article))
	a.db.AddArticle(*article)
}
