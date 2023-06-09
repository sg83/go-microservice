// Code generated by mockery v2.23.2. DO NOT EDIT.

package mocks

import (
	"github.com/sg83/go-microservice/article-api/data"
	mock "github.com/stretchr/testify/mock"
)

// ArticlesData is an autogenerated mock type for the ArticlesData type
type ArticlesData struct {
	mock.Mock
}

// AddArticle provides a mock function with given fields: ar
func (_m *ArticlesData) AddArticle(ar data.Article) error {
	ret := _m.Called(ar)

	var r0 error
	if rf, ok := ret.Get(0).(func(data.Article) error); ok {
		r0 = rf(ar)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *ArticlesData) Close() {
	_m.Called()
}

// GetArticleByID provides a mock function with given fields: id
func (_m *ArticlesData) GetArticleByID(id int) (*data.Article, error) {
	ret := _m.Called(id)

	var r0 *data.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*data.Article, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *data.Article); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Article)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetArticlesForTagAndDate provides a mock function with given fields: tag, date
func (_m *ArticlesData) GetArticlesForTagAndDate(tag string, date string) ([]int, error) {
	ret := _m.Called(tag, date)

	var r0 []int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]int, error)); ok {
		return rf(tag, date)
	}
	if rf, ok := ret.Get(0).(func(string, string) []int); ok {
		r0 = rf(tag, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(tag, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelatedTagsForTag provides a mock function with given fields: tag, articles
func (_m *ArticlesData) GetRelatedTagsForTag(tag string, articles []int) ([]string, error) {
	ret := _m.Called(tag, articles)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []int) ([]string, error)); ok {
		return rf(tag, articles)
	}
	if rf, ok := ret.Get(0).(func(string, []int) []string); ok {
		r0 = rf(tag, articles)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []int) error); ok {
		r1 = rf(tag, articles)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewArticlesData interface {
	mock.TestingT
	Cleanup(func())
}

// NewArticlesData creates a new instance of ArticlesData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewArticlesData(t mockConstructorTestingTNewArticlesData) *ArticlesData {
	mock := &ArticlesData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
