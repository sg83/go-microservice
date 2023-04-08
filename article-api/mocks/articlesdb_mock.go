package mocks

import (
	"github.com/sg83/go-microservice/article-api/models"
	"github.com/stretchr/testify/mock"
)

// mockArticlesDb is a mock implementation of the database interface
type ArticlesDbMock struct {
	mock.Mock
}

// Mocked method for getting article with given id
func (m *ArticlesDbMock) GetArticleByID(id int) (*models.Article, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *ArticlesDbMock) AddArticle(ar models.Article) error {
	args := m.Called(ar)
	return args.Error(0)
}

func (m *ArticlesDbMock) init() error {
	args := m.Called()
	return args.Error(0)
}
func (m *ArticlesDbMock) Close() {
	m.Called()
}
