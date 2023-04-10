package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/data"
	"github.com/sg83/go-microservice/article-api/utils"
	"go.uber.org/zap"
)

// KeyArticle is a key used for the Article object in the context
type KeyArticle struct{}

type Articles struct {
	l  *zap.Logger
	db data.ArticlesData
	v  *data.Validation
}

func NewArticles(l *zap.Logger, db data.ArticlesData, v *data.Validation) *Articles {
	return &Articles{l, db, v}
}

// Get retrieves an article by ID.
//
// swagger:operation GET /articles/{id} articles Get
//
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the article to retrieve
//     required: true
//     type: integer
//
// responses:
//
//	'200':
//	  description: Article retrieved successfully
//	  schema:
//	    "$ref": "#/definitions/Article"
//	'404':
//	  description: Article not found
//	  schema:
//	    "$ref": "#/definitions/GenericError"
//	'500':
//	  description: Internal server error
//	  schema:
//	    "$ref": "#/definitions/GenericError"
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
		utils.ToJSON(&utils.GenericError{Message: err.Error()}, w)
		return
	}

	err = utils.ToJSON(article, w)
	if err != nil {
		// we should never be here but log the error just incase
		a.l.Error("Unable to serialize product", zap.String(" Error: ", err.Error()))
		return
	}
}

// Create adds a new article.
//
// swagger:operation POST /articles articles Create
//
// ---
// parameters:
//   - name: article
//     in: body
//     description: Article to create
//     required: true
//     schema:
//     "$ref": "#/definitions/Article"
//
// responses:
//
//	'200':
//	  description: Article created successfully
//	'400':
//	  description: Invalid request payload
//	  schema:
//	    "$ref": "#/definitions/GenericError"
//	'500':
//	  description: Internal server error
//	  schema:
//	    "$ref": "#/definitions/GenericError"
func (a *Articles) Create(w http.ResponseWriter, r *http.Request) {

	a.l.Info("Create article", zap.Any("article:", r.Context().Value(KeyArticle{})))

	// fetch the article from the context
	article, ok := r.Context().Value(KeyArticle{}).(*data.Article)
	if !ok {
		// handle the case where the value is not of the expected type
		a.l.Error("Error fetching object from context")
		return
	}

	a.l.Info("Inserting ", zap.Any("article: ", article))
	a.db.AddArticle(*article)
}
