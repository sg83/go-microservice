package handlers

import (
	"context"
	"net/http"

	"github.com/sg83/go-microservice/article-api/data"
	"github.com/sg83/go-microservice/article-api/utils"
	"go.uber.org/zap"
)

// MiddlewareValidateProduct validates the article in the request and calls next if ok
func (a *Articles) MiddlewareValidateArticle(next http.Handler) http.Handler {
	a.l.Info("Validating article")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		article := &data.Article{}

		err := utils.FromJSON(article, r.Body)
		if err != nil {
			a.l.Error("Deserializing article ", zap.String("Error: ", err.Error()))
			rw.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&utils.GenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := a.v.Validate(article)
		if len(errs) != 0 {
			a.l.Error("Validating article", zap.Strings("Errors: ", errs.Errors()))

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			utils.ToJSON(&utils.ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyArticle{}, article)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
