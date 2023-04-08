package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/database"
	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

func (a *Articles) GetTagSummary(w http.ResponseWriter, r *http.Request) {
	a.l.Info("Get tag summary")
	vars := mux.Vars(r)
	tag := vars["tag"]
	dateStr := vars["date"]

	a.l.Info("Get tag summary", zap.String("tag:", tag), zap.String("date:", dateStr))

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	articlesIds, err := a.db.GetArticlesForTagAndDate(tag, date)
	//article, err := a.db.GetArticlesForTagAndDate(tagName, date)
	if err != nil {
		http.Error(w, "Articles with given tag not found", http.StatusNotFound)
		w.WriteHeader(http.StatusInternalServerError)
		database.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	a.l.Info("Get tag summary", zap.Any("Articles with tag:", articlesIds))

	relatedTags, err := a.db.GetRelatedTagsForTagAndDate(tag, date)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Related tags not found", http.StatusNotFound)
		w.WriteHeader(http.StatusInternalServerError)
		database.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	a.l.Info("Get tag summary", zap.Any("Related tags:", relatedTags))

	tagSummary := models.Tag{
		Tag:         tag,
		Count:       len(articlesIds),
		Articles:    articlesIds,
		RelatedTags: relatedTags,
	}

	w.Header().Set("Content-Type", "application/json")
	err = database.ToJSON(tagSummary, w)
	if err != nil {
		// we should never be here but log the error just incase
		a.l.Error("Unable to serialize tagSummary", zap.String(" Error: ", err.Error()))
		return
	}
}

/*
func getArticlesForTagAndDate(db *Articles, tagName string, date time.Time) ([]Article, error) {
	rows, err := db.Query("SELECT id, title, date, body, tags FROM articles WHERE $1 = ANY(tags) AND date = $2", tagName, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Date, &article.Body, pq.Array(&article.Tags))
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
*/
