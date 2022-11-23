package stats

import (
	"testing"
)

func TestMedianAgeExplict(t *testing.T) {

	testCases := []struct {
		name         string
		getAges      func() []int
		isActual     bool
		expectNoData bool
		expect       float32
	}{
		{
			name: "JustOne",
			getAges: func() []int {
				return []int{50}
			},
			expectNoData: false,
			isActual:     true,
			expect:       float32(50),
		},
		{
			name: "OddCount",
			getAges: func() []int {
				return []int{46, 4, 7, 11, 33}
			},
			expectNoData: false,
			isActual:     true,
			expect:       float32(11),
		},
		{
			name: "EvenCountNotActual",
			getAges: func() []int {
				return []int{46, 4, 7, 11, 33, 12}
			},
			expectNoData: false,
			isActual:     false,
			expect:       float32(11.5),
		},
		{
			name: "EvenCountActual",
			getAges: func() []int {
				return []int{46, 33, 14, 111, 33, 7}
			},
			expectNoData: false,
			isActual:     true,
			expect:       float32(33),
		},
		{
			name: "Int1To100",
			getAges: func() []int {
				res := make([]int, 100)
				for i := 1; i <= 100; i++ {
					res[i-1] = i
				}
				return res
			},
			expectNoData: false,
			isActual:     false,
			expect:       float32(50.5),
		},
		{
			name: "Int1To100Except50",
			getAges: func() []int {
				res := make([]int, 100-1)
				for i := 1; i <= 100; i++ {
					if i < 50 {
						res[i-1] = i
					} else if i > 50 {
						res[i-2] = i
					}
				}
				return res
			},
			expectNoData: false,
			isActual:     true,
			expect:       float32(51),
		},
		{
			name: "Int1To99",
			getAges: func() []int {
				res := make([]int, 99)
				for i := 1; i <= 99; i++ {
					res[i-1] = i
				}
				return res
			},
			expectNoData: false,
			isActual:     true,
			expect:       float32(50),
		},
		{
			name: "Int100To125",
			getAges: func() []int {
				res := make([]int, 26)
				for i := 100; i <= 125; i++ {
					res[i-100] = i
				}
				return res
			},
			expectNoData: false,
			isActual:     false,
			expect:       float32(112.5),
		},
		{
			name: "ExpectNoData",
			getAges: func() []int {
				return []int{}
			},
			expectNoData: true,
			isActual:     false,
			expect:       float32(0),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ages := tc.getAges()
			// log.Println("ages before:", ages)
			shuffle(ages)
			// log.Println("ages after shuffle:", ages)

			medianAge := NewMedianAge()
			for _, n := range ages {
				medianAge.Add(n)
			}

			got, isActual, err := medianAge.Calc()
			if err != nil {
				if err == ErrNoData {
					if !tc.expectNoData {
						t.Errorf("expect data, got ErrNoData")
					}
				} else {
					t.Errorf("unexpected error from MedianAge.Calc(): %v", err)
				}
			}

			if isActual != tc.isActual {
				t.Errorf("isActual: expect %v, got %v", tc.isActual, isActual)
			}

			if float32(got) != tc.expect {
				t.Errorf("median: expect %.2f, got %.2f", tc.expect, got)
			}
		})
	}

}
