package handlers_models

import "github.com/getnexar/golang-programming-task/doc-index/pkg/index"

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type SearchResponse struct {
	Results []index.IndexedDocument `json:"results"`
}
