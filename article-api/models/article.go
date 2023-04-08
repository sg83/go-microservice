package models

type Article struct {
	// Unique identifier for the product
	//
	// required: true
	// min: 1
	ID int `json:"id"`

	// the title for this article
	//
	// required: true
	// max length: 500
	Title string `json:"title" validate:"required"`

	// the date of the article
	//
	// required: true
	Date string `json:"date"`

	// the body for this article
	//
	// required: true
	// max length: 10000
	Body string `json:"body" validate:"required"`

	// the tags for the article
	//
	// required: false
	Tags []string `json:"tags"`
}
