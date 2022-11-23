package models

import "time"

type RunStats struct {
	InvalidUrls       []string      `json:"invalid_urls"`
	DuplicateUrls     []string      `json:"duplicate_urls"`
	DuplicateUrlCount int           `json:"duplicate_url_count"`
	ReadTimeCost      time.Duration `json:"read_time_cost"`
	NumThreads        int           `json:"num_threads"`
}
