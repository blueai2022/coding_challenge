package models

type DataDigest struct {
	TotalAgeCounts int64       `json:"total_age_counts"`
	AgeCounts      [201]int64  `json:"age_counts"`      //0 to 200
	AgePersonName  [201]string `json:"age_person_name"` //0 to 200
}
