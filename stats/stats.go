package stats

import "errors"

var (
	ErrNoDataForAverage = errors.New("cannot calculate average for 0 population")
)

func AverageAge(ageCounts [201]int64) (float32, error) {
	var ageTotal int64
	var totalCount int64

	for age, ageCount := range ageCounts {
		totalCount += ageCount
		ageTotal += int64(age) * ageCount
	}

	if totalCount == 0 {
		return float32(-1), ErrNoDataForAverage
	}

	return float32(ageTotal / totalCount), nil
}
