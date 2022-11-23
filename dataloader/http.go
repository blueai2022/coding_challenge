package dataloader

import (
	"net/http"
	"time"

	"github.com/blucv2022/crowdstats/models"
)

const (
	timeoutMins = 10
	CSVMimeType = "text/csv"
)

type HttpLoader struct {
}

func NewHttpLoader() Loader {
	return &HttpLoader{}
}

func (loader *HttpLoader) LoadDigest(url string) (*models.DataDigest, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", CSVMimeType)

	client := http.Client{
		Timeout: timeoutMins * time.Minute,
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		if rsp.StatusCode == http.StatusNotFound {
			return nil, &LoadError{Err: ErrFileNotFound}
		}
		return nil, &LoadError{Err: ErrFileLoadFailure}
	}

	digest, err := Digest(rsp.Body)
	if err != nil {
		return nil, &LoadError{Err: err}
	}

	return digest, nil
}
