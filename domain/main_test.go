package domain

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/blucv2022/crowdstats/dataloader"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, https bool) *httptest.Server {
	mux := http.NewServeMux()

	// to test any HTTP error code specified in URL
	mux.HandleFunc("/TESTERR/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			t.Fatal("test error code not found in bad url")
		}

		code, err := strconv.Atoi(parts[2])
		if err != nil {
			t.Fatal("cannot parse test error code from bad url", err)
		}

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]string{"message": "testing http server errors"}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			t.Fatal("cannot marshal json to response:", err)
		}
		w.Write(jsonResp)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != dataloader.CSVMimeType {
			t.Errorf("Expected Accept: %s, got: %s", dataloader.CSVMimeType, r.Header.Get("Accept"))
		}

		filePath := "../" + r.URL.Path
		content, err := readTestData(t, filePath)
		if err != nil {
			if err == os.ErrNotExist {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(content)
		if err != nil {
			t.Error("cannot write response:", err)
		}
	})

	if https {
		//TODO: configure a widely accepted SSL cert to get around below issue on dev Mac:
		//  x509: “Acme Co” certificate is not trusted
		return httptest.NewTLSServer(mux)
	} else {
		return httptest.NewServer(mux)
	}
}

func getDataSrcUrls(t *testing.T, filePath string) []string {
	// open source url file
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal("cannot open source url file:", err)
	}
	defer f.Close()

	// read source urls line by line
	scanner := bufio.NewScanner(f)
	res := []string{}
	for scanner.Scan() {
		//trim and length check to protect against blank lines in file,
		//especially ones at the end
		text := strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			res = append(res, text)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal("cannot read source url lines:", err)
	}

	return res
}
