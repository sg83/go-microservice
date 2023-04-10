package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sg83/go-microservice/article-api/data"
	"github.com/sg83/go-microservice/article-api/mocks"
	"go.uber.org/zap"
)

func TestGetTagSummary(t *testing.T) {

	tt := []struct {
		tagSummary *data.Tag
		date       string
		status     int
		err        string
	}{
		{
			tagSummary: &data.Tag{
				Tag:         "health",
				Count:       3,
				Articles:    []int{1, 3, 5},
				RelatedTags: []string{"yoga", "fitness"},
			},
			date:   "20220512",
			status: 200,
			err:    "",
		},
		{
			tagSummary: &data.Tag{
				Tag:         "lifestyle",
				Count:       2,
				Articles:    []int{1, 4},
				RelatedTags: []string{"yoga", "fitness"},
			},
			date:   "20220512",
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
		req, err := http.NewRequest("GET", "/tags/?v="+tc.tagSummary.Tag+"/?v="+tc.date, nil)
		if err != nil {
			t.Fatal(err)
		}

		// create a mock Articles struct with a mock database interface
		mockdb := new(mocks.ArticlesData)
		mockdb.On("GetArticlesForTagAndDate", tc.tagSummary.Tag, tc.date).Return(tc.tagSummary.Articles, nil)
		mockdb.On("GetRelatedTagsForTag", tc.tagSummary.Tag, tc.tagSummary.Articles).Return(tc.tagSummary.RelatedTags, nil)

		articles := &Articles{logger, mockdb, nil}

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"tag":  tc.tagSummary.Tag,
			"date": tc.date,
		}
		req = mux.SetURLVars(req, vars)

		// call the Get function with the mock response writer and request
		articles.GetTagSummary(w, req)

		// check that the response status code is 200 OK
		if w.Code != tc.status {
			t.Errorf("Expected status code %d but got %d", tc.status, w.Code)
		}

		// check that the response body contains the expected article
		expected := tc.tagSummary
		actual := &data.Tag{}
		err = json.NewDecoder(w.Body).Decode(actual)
		if err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected tag Summary %v but got %v", expected, actual)
		}
		t.Logf("Test passed for tag %s", tc.tagSummary.Tag)
	}
}
