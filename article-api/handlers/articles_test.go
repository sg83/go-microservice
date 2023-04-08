package handlers

/*
func TestGetArticle(t *testing.T) {

	tt := []struct {
		id     string
		title  string
		body   string
		date   string
		tags   []string
		status int
		err    string
	}{
		{
			id:     "1",
			title:  "Article1",
			body:   "This article is about health and fitness.",
			date:   "20-02-2023",
			tags:   []string{"health", "fitness"},
			status: 200,
			err:    "",
		},
	}

	// create a mock response writer
	w := httptest.NewRecorder()

	// create a mock request with a URL containing an article ID
	req, err := http.NewRequest("GET", "/articles/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a mock Articles struct with a mock database interface
	mockdb := new(mocks.ArticlesDbMock)
	mockdb.On("GetArticleByID", 1).Return(tt[1])
	articles := &Articles{nil, mockdb}

	// call the Get function with the mock response writer and request
	articles.Get(w, req)

	// check that the response status code is 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

	// check that the response body contains the expected article
	expected := &models.Article{ID: 1, Title: "Test Article", Body: "Lorem ipsum dolor sit amet."}
	actual := &models.Article{}
	err = json.NewDecoder(w.Body).Decode(actual)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected article %v but got %v", expected, actual)
	}
}
*/
/*
func TestGet(t *testing.T) {
	tt := []struct {
		id     string
		title  string
		body   string
		date   string
		tags   []string
		status int
		err    string
	}{
		{
			id:     "1",
			title:  "Article1",
			body:   "This article is about health and fitness.",
			date:   "20-02-2023",
			tags:   []string{"health", "fitness"},
			status: 200,
			err:    "",
		},
	}

	for _, tc := range tt {
		req, err := http.NewRequest("GET", "localhost:8080/articles?v="+tc.id, nil)
		if err != nil {
			t.Fatal("Could not create request: %v", err)
		}

		rec := httptest.NewRecorder()
		Get(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if tc.err != "" {
			//do something
			return
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", res)
		}

		a, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Could not read response %v", err)
		}
		t.Log("output=%v", a)
	}
}

func TestRouting(t *testing.T) {

}
*/
// mockDatabase is a mock implementation of the database interface
