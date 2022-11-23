package domain

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blucv2022/crowdstats/dataloader"
)

const (
	testSvrAddrPlaceHolder    = "<test-server-addr>"
	testFilePathRoot          = "../testdata"
	generatedTestFilePathRoot = "../testdata/generated"
)

type testCase struct {
	name                string
	generatedTestFiles  bool
	srcUrlFileFunc      func() (string, bool, []string)
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

func TestStaticFileSets(t *testing.T) {
	server := newTestServer(t, false)
	if server == nil {
		t.Fatal("cannot start test http server")
	}
	defer server.Close()

	testCases := []testCase{
		{
			name:               "OriginalSetTestFiles",
			generatedTestFiles: false,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				filePath = fmt.Sprintf("%s/http_src_urls.txt", testFilePathRoot)
				isTempFile = false
				return
			},
			expectNoData:       false,
			expInvalidUrlCount: 2,
			expInvalidUrls: []string{
				"http://<test-server-addr>/testdata/files/file6_bad.csv",
				"http://<test-server-addr>/testdata/files/file9_bad.csv",
			},
			expDupeUrlCount:    0,
			expMedianAge:       float32(31),
			expAverageAge:      float32(33),
			expMedianAgePerson: true,
		},
		{
			name:               "DupeUrls",
			generatedTestFiles: false,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				filePath = fmt.Sprintf("%s/http_src_urls_dupe.txt", testFilePathRoot)
				isTempFile = false
				return
			},
			expectNoData:       false,
			expInvalidUrlCount: 0,
			expDupeUrlCount:    4,
			expDupeUrls: []string{
				"http://<test-server-addr>/testdata/files/file1.csv",
				"http://<test-server-addr>/testdata/files/file1.csv",
				"http://<test-server-addr>/testdata/files/file5.csv",
				"http://<test-server-addr>/testdata/files/file5.csv",
			},
			expMedianAge:       float32(31),
			expAverageAge:      float32(33),
			expMedianAgePerson: true,
		},
		{
			name:               "HTTPErrors",
			generatedTestFiles: false,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				filePath = fmt.Sprintf("%s/http_src_urls_bad.txt", testFilePathRoot)
				isTempFile = false
				return
			},
			expectNoData:       false,
			expInvalidUrlCount: 8,
			expInvalidUrls: []string{
				"http://<test-server-addr>/TESTERR/401/files/file3.csv",
				"http://<test-server-addr>/TESTERR/403/files/file2.csv",
				"http://<test-server-addr>/TESTERR/404/files/file2.csv",
				"http://<test-server-addr>/TESTERR/500/files/file3.csv",
				"http://<test-server-addr>/TESTERR/501/files/file3.csv",
				"http://<test-server-addr>/TESTERR/502/files/file3.csv",
				"http://<test-server-addr>/TESTERR/503/files/file3.csv",
				"http://<test-server-addr>/TESTERR/504/files/file3.csv",
			},
			expDupeUrlCount:     0,
			expMedianAge:        float32(90),
			expAverageAge:       float32(89),
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{"Julia DAVILA"},
		},
	}

	//Before running test cases, clean up all generated files under path root for a fresh run
	cleanupTempDir(t, generatedTestFilePathRoot)

	for _, tc := range testCases {
		runTestCase(t, tc, server)
	}
}

func TestInvalidUrlsAndPersonName(t *testing.T) {
	server := newTestServer(t, false)
	if server == nil {
		t.Fatal("cannot start test http server")
	}
	defer server.Close()

	testCases := []testCase{
		{
			name:               "2InvalidFiles",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pike",25
Ken,Thompson,210
"Robert","Griesemer",50
`,
					`fname,lname,age
"Koko","Pike",-25
Kenny,Thompson,33
"Roberta","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 2,
			expDupeUrlCount:    0,
		},
		{
			name:               "2Good1Bad",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Bob","Kick",4
Henry,Thompson,30
"Robert","Brady",50
`,
					`fname,lname,age
"Rob","Pike",15
Ken,Thompson,31
"Robert","Griesemer",57
`,
					`fname,lname,age
"Koko","Pike",-25
Kenny,Thompson,33
"Roberta","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       false,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
			expMedianAge:       30.5,
			expAverageAge:      31,
		},
		{
			name:               "MatchMedianAgeNamesOneOf",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Bob","Kick",4
Henry,Thompson,30
"Robert","Brady",50
`,
					`fname,lname,age
"Koko","Pike", 31
Kenny,Thompson,33
Alice,Wonders, 31
"Roberta","Griesemer",50
`,
					`fname,lname,age
"Rob","Pike",15
Ken,Thompson,31
"Robert","Griesemer",57
`}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:        false,
			expInvalidUrlCount:  0,
			expDupeUrlCount:     0,
			expMedianAge:        31,
			expAverageAge:       33,
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{"Ken Thompson", "Alice Wonders", "Koko Pike"},
		},
	}

	//Before running test cases, clean up all generated files under path root for a fresh run
	cleanupTempDir(t, generatedTestFilePathRoot)

	for _, tc := range testCases {
		runTestCase(t, tc, server)
	}
}

func TestDataValidation(t *testing.T) {
	server := newTestServer(t, false)
	if server == nil {
		t.Fatal("cannot start test http server")
	}
	defer server.Close()

	testCases := []testCase{
		{
			name:               "InvalidAgesNegative",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pike",25
Ken,Thompson,-10
"Robert","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
		},
		{
			name:               "InvalidAgesOver200",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pike",25
Ken,Thompson,210
"Robert","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
		},
		{
			name:               "InvalidPersonNameWithSpecialCharsEx1",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pi**ke",25
Ken,Thompson,21
"Robert","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
		},
		{
			name:               "InvalidPersonNameWithSpecialCharsEx2",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pike",25
Ken,Thomp$$son,21
"Robert","Griesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
		},
		{
			name:               "InvalidPersonNameWithSpecialCharsEx3",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Pike",25
Ken,Thompson,21
"Robert","((Griese))mer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:       true,
			expInvalidUrlCount: 1,
			expDupeUrlCount:    0,
		},
		{
			name:               "ValidPersonNameWithSpecialChars",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Rob","Downey Jr.",25
Ken,Thomp-son,21
"Robert","G'riesemer",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:        false,
			expInvalidUrlCount:  0,
			expDupeUrlCount:     0,
			expMedianAge:        25,
			expAverageAge:       32,
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{"Rob Downey Jr."},
		},
		{
			name:               "ValidPersonNameInUnicode",
			generatedTestFiles: true,
			srcUrlFileFunc: func() (filePath string, isTempFile bool, tmpDataFilePaths []string) {
				data := []string{
					`fname,lname,age
"Béa","BURNS",25
张秀英, SCOTT,30
"English","Normal",50
`,
				}

				filePath, tmpDataFilePaths = generateTestFiles(t, false, data)
				isTempFile = true
				return
			},
			expectNoData:        false,
			expInvalidUrlCount:  0,
			expDupeUrlCount:     0,
			expMedianAge:        30,
			expAverageAge:       35,
			expMedianAgePerson:  true,
			expMedianAgePersons: []string{"张秀英 SCOTT"},
		},
	}

	//Before running test cases, clean up all generated files under path root for a fresh run
	cleanupTempDir(t, generatedTestFilePathRoot)

	for _, tc := range testCases {
		runTestCase(t, tc, server)
	}
}

func runTestCase(t *testing.T, tc testCase, server *httptest.Server) {
	t.Run(tc.name, func(t *testing.T) {
		if !tc.generatedTestFiles &&
			tc.expInvalidUrlCount != len(tc.expInvalidUrls) {
			t.Fatal("invalid test setup: expInvalidUrlCount and expInvalidUrls have conflicts")
		}

		if !tc.generatedTestFiles &&
			tc.expDupeUrlCount != len(tc.expDupeUrls) {
			t.Fatal("invalid test setup: expDupeUrlCount and expDupeUrls have conflicts")
		}

		// Okay to expect median person (true) and not specify list of possible names (too many to list)
		// However, not the other way around
		if !tc.expMedianAgePerson && len(tc.expMedianAgePersons) > 0 {
			t.Fatal("invalid test setup: expMedianAgePerson and expMedianAgePersons have conflicts")
		}

		srcUrlFile, isTempFile, tmpDataFilePaths := tc.srcUrlFileFunc()
		srcUrls := getDataSrcUrls(t, srcUrlFile)
		if isTempFile {
			// cleanup tmp test files for this test only; to allow for parallel tests
			defer cleanupTempFiles(t, []string{srcUrlFile})
			defer cleanupTempFiles(t, tmpDataFilePaths)
		}

		// replace Server Address place holders with actual httptest.Server URL, dynamic at runtime
		testSrcUrls := make([]string, len(srcUrls))
		testServerAddr := server.Listener.Addr().String()
		for i := range srcUrls {
			testSrcUrls[i] = strings.Replace(srcUrls[i],
				testSvrAddrPlaceHolder,
				testServerAddr,
				1)
		}

		loader := dataloader.NewHttpLoader()

		summarizer := NewSummarizer(loader)
		dataSum, runStats := summarizer.Run(testSrcUrls)

		//reverse Server Address replacement to show correct invalid source urls
		for i := range runStats.InvalidUrls {
			runStats.InvalidUrls[i] =
				strings.Replace(runStats.InvalidUrls[i],
					testServerAddr,
					testSvrAddrPlaceHolder,
					1,
				)
		}

		//reverse Server Address replacement to show correct dupe source urls
		for i := range runStats.DuplicateUrls {
			runStats.DuplicateUrls[i] =
				strings.Replace(runStats.DuplicateUrls[i],
					testServerAddr,
					testSvrAddrPlaceHolder,
					1,
				)
		}

		if runStats == nil {
			t.Error("Expected valued run stats, got nil")
		} else {
			if len(runStats.InvalidUrls) != tc.expInvalidUrlCount {
				t.Errorf("Expected count of %d invalid urls, got %d",
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

		if dataSum == nil {
			t.Fatal("Expected valued data summary, got nil")
		}

		// check no data
		if tc.expectNoData {
			if tc.expMedianAgePerson {
				t.Fatal("Incorrect setup: expectNoData and expMedianAgePerson are both true")
			}

			if (!dataSum.IsNA) ||
				dataSum.MedianAge > 0.0 ||
				dataSum.AverageAge > 0.0 ||
				len(dataSum.MedianAgePerson) > 0 ||
				dataSum.IsMedianAgeActual {

				t.Errorf("Expected no data, got data %v", dataSum)
			}
			return
		}

		if dataSum.IsNA != tc.expectNoData {
			t.Errorf("IsNA expected %t, got %t",
				tc.expectNoData,
				dataSum.IsNA,
			)
		}

		if dataSum.IsMedianAgeActual != tc.expMedianAgePerson {
			t.Errorf("IsMedianAgeActual expected %t, got %t",
				tc.expMedianAgePerson,
				dataSum.IsMedianAgeActual,
			)
		}

		if dataSum.MedianAge != tc.expMedianAge {
			t.Errorf("Median Age expected %.2f, got %.2f",
				tc.expMedianAge,
				dataSum.MedianAge,
			)
		}

		if dataSum.AverageAge != tc.expAverageAge {
			t.Errorf("Average Age expected %.2f, got %.2f",
				tc.expAverageAge,
				dataSum.AverageAge,
			)
		}

		if !tc.expMedianAgePerson && len(dataSum.MedianAgePerson) > 0 {
			t.Errorf("Person with Median Age expected blank, got %s",
				dataSum.MedianAgePerson,
			)
		}

		if tc.expMedianAgePerson && len(dataSum.MedianAgePerson) == 0 {
			t.Error("Person with Median Age expected to be valued, got blank")
		}

		if len(tc.expMedianAgePersons) > 0 &&
			!contains(tc.expMedianAgePersons, dataSum.MedianAgePerson) {
			t.Errorf("Person with Median Age expected to be one of %v, got %s",
				tc.expMedianAgePersons,
				dataSum.MedianAgePerson,
			)
		}
	})

}
