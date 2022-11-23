package models

type Summary struct {
	IsNA              bool    `json:"is_na"`
	AverageAge        float32 `json:"average_age"`
	MedianAge         float32 `json:"median_age"`
	IsMedianAgeActual bool    `json:"is_median_age_actual"`
	MedianAgePerson   string  `json:"median_age_person"`
}
