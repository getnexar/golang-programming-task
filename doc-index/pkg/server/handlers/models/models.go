package handlers_models

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type SearchResponse struct {
	Results [][2]string `json:"results"`
}
