package stats

import (
	"errors"
	"log"

	"github.com/blucv2022/crowdstats/val"
)

var (
	ErrNoData = errors.New("no age data")
)

// internal errors
var (
	errQueryOutOfBounds = errors.New("index out of total count range")
)

// range-bound median implementation for Age
type MedianAge struct {
	ageCounts  [201]int64 //0 to 200
	totalCount int64
}

func NewMedianAge() *MedianAge {
	return &MedianAge{}
}

func (mdn *MedianAge) AddAll(ageCounts [201]int64, totalCount int64) {
	mdn.ageCounts = ageCounts
	mdn.totalCount = totalCount
}

func (mdn *MedianAge) Add(age int) {
	if err := val.ValidateAge(age); err != nil {
		log.Fatalf("invalid age: %d", age)
	}

	mdn.ageCounts[age]++
	mdn.totalCount++
}

func (mdn *MedianAge) Calc() (median float64, isActual bool, err error) {
	//edge case: no age data
	if mdn.totalCount == 0 {
		err = ErrNoData
		return
	}

	if mdn.totalCount%2 == 0 {
		midNumIdx1 := mdn.totalCount/2 - 1

		mid2Ages, err := mdn.getTwoConsecVals(midNumIdx1)
		if err != nil {
			log.Fatal("cannot get median age:", err)
		}

		isActual = mid2Ages[0] == mid2Ages[1]

		if isActual {
			median = float64(mid2Ages[0])
		} else {
			median = float64(mid2Ages[0]+mid2Ages[1]) / 2
		}
	} else {
		//odd count, median is an actual value
		isActual = true

		midNumIdx := (mdn.totalCount+1)/2 - 1

		midAges, err := mdn.getTwoConsecVals(midNumIdx)
		if err != nil {
			log.Fatal("cannot get median age:", err)
		}
		median = float64(midAges[0])
	}

	return
}

func (mdn *MedianAge) getTwoConsecVals(startIdx int64) (*[2]int, error) {
	//check out of bounds: only for first number at pos idx
	//2nd number is optional; so when totalCount = 1 and idx = 0, it is okay
	if startIdx+1 > mdn.totalCount {
		return nil, errQueryOutOfBounds
	}

	//init to 2 invalid values
	resAges := &[2]int{-1, -1}

	var runningTotal int64
	for i := 0; i < len(mdn.ageCounts); i++ {
		//first age number not found
		if resAges[0] < 0 {
			runningTotal += mdn.ageCounts[i]
			if runningTotal >= startIdx+1 {
				resAges[0] = i

				//if 2nd age number is also in range with runningTotal
				// then 2nd age number = 1st age number
				if runningTotal >= (startIdx+1)+1 {
					resAges[1] = i
					break
				}
			}
		} else {
			//first age number previously found, look for second
			// next age number with valid count -> 2nd age number
			if mdn.ageCounts[i] > 0 {
				resAges[1] = i
				break
			}
		}
	}

	return resAges, nil
}
