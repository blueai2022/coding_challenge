package domain

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func cleanupTempFiles(t *testing.T, filePaths []string) {
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			t.Error("cannot remove temp test file:", err)
		}
	}
}

func cleanupTempDir(t *testing.T, rootPath string) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		// fail test: since may affect new testing results
		t.Fatal("cannot read test temp dir:", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// fail test: since does not know about sub-dir being used;
			// may affect new testing results
			t.Fatal("unexpected sub-directory under test temp dir")
		} else {
			filePath := rootPath + "/" + entry.Name()
			//warn about previous tests not cleaning up the files after done
			//so that main test code gets improved
			t.Log("warn: uncleaned temp file -", filePath)

			err := os.Remove(filePath)
			if err != nil {
				t.Error("cannot remove temp source urls file:", err)
			}
		}
	}
}

func writeTestFile(t *testing.T, filePathString, data string) {
	err := os.WriteFile(filePathString, []byte(data), 0644)
	if err != nil {
		t.Fatal("cannot write test file:", err)
	}
}

func readTestData(t *testing.T, filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		//file found, cannot be read, error
		t.Error("cannot read file:", err)
	}
	return content, nil
}

func generateTestFiles(t *testing.T, https bool, data []string) (string, []string) {
	if len(data) == 0 {
		t.Fatal("test data not provided")
	}

	testId := randomString(8)
	srcUrlFile := fmt.Sprintf("%s/src_urls_%s.txt",
		generatedTestFilePathRoot,
		testId,
	)

	protocol := "http"
	if https {
		protocol = "https"
	}

	var srcUrls string
	dataFilePaths := make([]string, len(data))

	for i, dt := range data {
		dataFilePath := fmt.Sprintf("%s/test_%s_file%d.csv",
			generatedTestFilePathRoot,
			testId,
			i+1,
		)

		//write data file
		writeTestFile(t, dataFilePath, dt)
		dataFilePaths[i] = dataFilePath

		//src urls, need to translate relative file path to url path by removing ../
		dataUrlPath := strings.Replace(dataFilePath, "../", "", 1)
		srcUrls += fmt.Sprintf("%s://%s/%s\n",
			protocol,
			testSvrAddrPlaceHolder,
			dataUrlPath,
		)
	}

	//write source urls file
	writeTestFile(t, srcUrlFile, srcUrls)

	return srcUrlFile, dataFilePaths
}
