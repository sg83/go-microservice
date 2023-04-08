package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/mocks"
	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

func TestGetArticle(t *testing.T) {

	tt := []struct {
		article *models.Article
		status  int
		err     string
	}{
		{
			article: &models.Article{
				ID:    1,
				Title: "Article1",
				Body:  "This article is about health and fitness.",
				Date:  "20-02-2023",
				Tags:  []string{"health", "fitness"},
			},
			status: 200,
			err:    "",
		},
		{
			article: &models.Article{
				ID:    2,
				Title: "Article2",
				Body:  "This article is about health and yoga.",
				Date:  "20-02-2023",
				Tags:  []string{"health", "yoga"},
			},
			status: 200,
			err:    "",
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	for _, tc := range tt {
		// create a mock response writer
		w := httptest.NewRecorder()

		// create a mock request with a URL containing an article ID
		req, err := http.NewRequest("GET", "/articles?v="+strconv.Itoa(tc.article.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		// create a mock Articles struct with a mock database interface
		mockdb := new(mocks.ArticlesData)
		mockdb.On("GetArticleByID", tc.article.ID).Return(tc.article, nil)
		articles := &Articles{logger, mockdb, nil}

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"id": strconv.Itoa(tc.article.ID),
		}
		req = mux.SetURLVars(req, vars)

		// call the Get function with the mock response writer and request
		articles.Get(w, req)

		// check that the response status code is 200 OK
		if w.Code != tc.status {
			t.Errorf("Expected status code %d but got %d", tc.status, w.Code)
		}

		// check that the response body contains the expected article
		expected := tc.article
		actual := &models.Article{}
		err = json.NewDecoder(w.Body).Decode(actual)
		if err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected article %v but got %v", expected, actual)
		}
		t.Logf("Test passed for article id %d", tc.article.ID)
	}
}
