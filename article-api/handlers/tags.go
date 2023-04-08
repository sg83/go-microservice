package handlers

import (
	"net/http"
	"regexp"
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

	re := regexp.MustCompile(`^(20[0-2][0-3]|1[2-9]|[2-9]\d)(\d{2})(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])$`)
	if !re.MatchString(dateStr) {
		a.l.Error("Date is not valid", zap.String("date:", dateStr))
		http.Error(w, "Date is not valid", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("20060102", dateStr)
	a.l.Info("Get tag summary", zap.String("date:", date.String()))
	if err != nil {
		a.l.Error("Could not parse date")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	articlesIds, err := a.db.GetArticlesForTagAndDate(tag, date)
	//article, err := a.db.GetArticlesForTagAndDate(tagName, date)
	if (err != nil) || (len(articlesIds) == 0) {
		a.l.Error("Articles with given tag not found")
		http.Error(w, "Articles with given tag not found", http.StatusNotFound)
		w.WriteHeader(http.StatusInternalServerError)
		database.ToJSON(&GenericError{Message: "Articles with given tag not found"}, w)
		return
	}
	a.l.Info("Get tag summary", zap.Any("Articles with tag:", articlesIds))

	relatedTags, err := a.db.GetRelatedTagsForTag(tag, articlesIds)
	if (err != nil) || (len(relatedTags) == 0) {
		a.l.Error("Related tags not found")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Related tags not found", http.StatusNotFound)
		w.WriteHeader(http.StatusInternalServerError)
		database.ToJSON(&GenericError{Message: "Related tags not found"}, w)
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
