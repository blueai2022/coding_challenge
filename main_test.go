package main

import (
	"fmt"
	"os"
	"sort"
	"testing"

	loadermock "github.com/blucv2022/crowdstats/dataloader/mock"
)

const (
	testFilePathRoot = "./testdata/"
)

type testCase struct {
	name                string
	srcUrlFile          string
	expectNoData        bool
	expInvalidUrlCount  int
	expInvalidUrls      []string
	expDupeUrlCount     int
	expDupeUrls         []string
	expMedianAge        float32
	expAverageAge       float32
	expMedianAgePerson  bool
	expMedianAgePersons []string
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestOriginalSet(t *testing.T) {
	tc := testCase{
		name:               "OriginalFileSet",
		srcUrlFile:         fmt.Sprintf("%s/file_src_urls.txt", testFilePathRoot),
		expectNoData:       false,
		expInvalidUrlCount: 2,
		expInvalidUrls: []string{
			"file://./testdata/files/file6_bad.csv",
			"file://./testdata/files/file9_bad.csv",
		},
		expDupeUrlCount:     0,
		expDupeUrls:         []string{},
		expMedianAge:        float32(31),
		expAverageAge:       float32(33),
		expMedianAgePerson:  true,
		expMedianAgePersons: []string{},
	}

	runTestCase(t, tc)
}

func TestFileLoaderBadAndDupeUrls(t *testing.T) {
	testCases := []testCase{
		{
			name:               "BadURLs",
			srcUrlFile:         fmt.Sprintf("%s/file_src_urls_bad.txt", testFilePathRoot),
			expectNoData:       false,
			expInvalidUrlCount: 5,
			expInvalidUrls: []string{
				"file://./testdata/files/NOTEXIST1.csv",
				"file://./testdata/files/NOTEXIST2.csv",
				"file://./testdata/files/NOTEXIST3.csv",
				"file://./testdata/files/NOTEXIST4.csv",
				"file://./testdata/files/NOTEXIST5.csv",
			},
			expDupeUrlCount:     0,
			expDupeUrls:         []string{},
			expMedianAge:        float32(90),
			expAverageAge:       float32(89),
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{"Julia DAVILA"},
		},
		{
			name:               "DupeURLs",
			srcUrlFile:         fmt.Sprintf("%s/file_src_urls_dupe.txt", testFilePathRoot),
			expectNoData:       false,
			expInvalidUrlCount: 0,
			expInvalidUrls:     []string{},
			expDupeUrlCount:    2,
			expDupeUrls: []string{
				"file://./testdata/files/file1.csv",
				"file://./testdata/files/file5.csv",
			},
			expMedianAge:        float32(31),
			expAverageAge:       float32(33),
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{},
		},
	}

	for _, tc := range testCases {
		runTestCase(t, tc)
	}
}

func runTestCase(t *testing.T, tc testCase) {
	t.Run(tc.name, func(t *testing.T) {
		if tc.expInvalidUrlCount != len(tc.expInvalidUrls) {
			t.Fatal("invalid test setup: expInvalidUrlCount and expInvalidUrls have conflicts")
		}

		if tc.expDupeUrlCount != len(tc.expDupeUrls) {
			t.Fatal("invalid test setup: expDupeUrlCount and expDupeUrls have conflicts")
		}

		// Okay to expect median person (true) and not specify list of possible names (too many to list)
		// However, not the other way around
		if !tc.expMedianAgePerson && len(tc.expMedianAgePersons) > 0 {
			t.Fatal("invalid test setup: expMedianAgePerson and expMedianAgePersons have conflicts")
		}

		loader := loadermock.NewFileLoader()
		summary, runStats := SummarizeAll(tc.srcUrlFile, loader)

		if runStats == nil {
			t.Error("Expected valued run stats, got nil")
		} else {
			if len(runStats.InvalidUrls) != tc.expInvalidUrlCount {
				t.Errorf("Expected %d invalid urls, got %d",
					tc.expInvalidUrlCount,
					len(runStats.InvalidUrls),
				)
			}

			if len(tc.expInvalidUrls) > 0 &&
				!equals(runStats.InvalidUrls, tc.expInvalidUrls, false) {
				t.Errorf("Expected invalid urls %v, got %v",
					tc.expInvalidUrls,
					runStats.InvalidUrls,
				)
			}

			if runStats.DuplicateUrlCount != tc.expDupeUrlCount {
				t.Errorf("Duplicate Urls count expected %d got %d",
					tc.expDupeUrlCount,
					runStats.DuplicateUrlCount,
				)
			}

			if len(tc.expDupeUrls) > 0 &&
				!equals(runStats.DuplicateUrls, tc.expDupeUrls, false) {
				t.Errorf("Expected duplicate urls %v, got %v",
					tc.expDupeUrls,
					runStats.DuplicateUrls,
				)
			}
		}

		if summary == nil {
			t.Fatal("Expected valued data summary, got nil")
		}

		// check no data
		if tc.expectNoData {
			if tc.expMedianAgePerson {
				t.Fatal("Incorrect setup: expectNoData and expMedianAgePerson are both true")
			}

			if (!summary.IsNA) ||
				summary.MedianAge > 0.0 ||
				summary.AverageAge > 0.0 ||
				len(summary.MedianAgePerson) > 0 ||
				summary.IsMedianAgeActual {

				t.Errorf("Expected no data, got data %v", summary)
			}
			return
		}

		if summary.IsNA != tc.expectNoData {
			t.Errorf("IsNA expected %t, got %t",
				tc.expectNoData,
				summary.IsNA,
			)
		}

		if summary.IsMedianAgeActual != tc.expMedianAgePerson {
			t.Errorf("IsMedianAgeActual expected %t, got %t",
				tc.expMedianAgePerson,
				summary.IsMedianAgeActual,
			)
		}

		if summary.MedianAge != tc.expMedianAge {
			t.Errorf("Median Age expected %.2f, got %.2f",
				tc.expMedianAge,
				summary.MedianAge,
			)
		}

		if summary.AverageAge != tc.expAverageAge {
			t.Errorf("Average Age expected %.2f, got %.2f",
				tc.expAverageAge,
				summary.AverageAge,
			)
		}

		if !tc.expMedianAgePerson && len(summary.MedianAgePerson) > 0 {
			t.Errorf("Person with Median Age expected blank, got %s",
				summary.MedianAgePerson,
			)
		}

		if tc.expMedianAgePerson && len(summary.MedianAgePerson) == 0 {
			t.Error("Person with Median Age expected to be valued, got blank")
		}

		if len(tc.expMedianAgePersons) > 0 &&
			!contains(tc.expMedianAgePersons, summary.MedianAgePerson) {
			t.Errorf("Person with Median Age expected to be one of %v, got %s",
				tc.expMedianAgePersons,
				summary.MedianAgePerson,
			)
		}
	})
}

func contains(strs []string, match string) bool {
	for _, s := range strs {
		if s == match {
			return true
		}
	}
	return false
}

func equals(a []string, b []string, ignoreOrder bool) bool {
	if len(a) != len(b) {
		return false
	}

	if ignoreOrder {
		sort.Strings(a)
		sort.Strings(b)
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
