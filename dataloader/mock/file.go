package loadermock

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/blucv2022/crowdstats/dataloader"
	"github.com/blucv2022/crowdstats/models"
)

type FileLoader struct {
}

func NewFileLoader() dataloader.Loader {
	return &FileLoader{}
}

func (loader *FileLoader) LoadDigest(fileUrl string) (*models.DataDigest, error) {
	// Problem: "file://..." urls do NOT support relative path
	// Normally just put full path urls in input data with all data source urls
	// However, for a take-home project, don't know reviewer's local file path
	// Workaround: convert relative path to absolute path url with help of working dir
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot get working dir:", err)
	}
	absPathUrl := strings.Replace(fileUrl, "//./", fmt.Sprintf("//%s/", workingDir), 1)

	parsedUrl, err := url.Parse(absPathUrl)
	if err != nil {
		return nil, &dataloader.LoadError{Err: err}
	}

	file, err := os.Open(parsedUrl.Path)
	if err != nil {
		return nil, &dataloader.LoadError{Err: err}
	}
	defer file.Close()

	digest, err := dataloader.Digest(file)
	if err != nil {
		return nil, &dataloader.LoadError{Err: err}
	}

	return digest, nil
}
