package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sg83/go-microservice/article-api/database"
	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

// MiddlewareValidateProduct validates the article in the request and calls next if ok
func (a *Articles) MiddlewareValidateArticle(next http.Handler) http.Handler {
	a.l.Info("Validating article")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		article := &models.Article{}

		fmt.Printf("req data %v\n", r.Body)
		err := database.FromJSON(article, r.Body)
		if err != nil {
			a.l.Error("Deserializing article ", zap.String("Error: ", err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
			database.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		fmt.Printf("req json data %v\n", article)
		// validate the product
		/*
			errs := a.v.Validate(article)
			if len(errs) != 0 {
				a.l.Error("Validating article", zap.Any("Errors: ", errs.Errors()))

				// return the validation messages as an array
				rw.WriteHeader(http.StatusUnprocessableEntity)
				database.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
				return
			}
		*/
		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyArticle{}, article)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
