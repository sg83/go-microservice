package data

type Tag struct {
	// Tag name
	Tag string `json:"tag"`
	// Number of articles having the tag for that day.
	Count int `json:"count"`
	// List of ids for the last 10 articles entered for that day.
	Articles []int `json:"articles"`
	// List of tags that are on the articles that the current tag is on for the same day.
	RelatedTags []string `json:"related_tags"`
}
