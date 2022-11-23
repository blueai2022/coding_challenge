package domain

import (
	"sync"

	"github.com/blucv2022/crowdstats/concurrency"
	"github.com/blucv2022/crowdstats/dataloader"
	"github.com/blucv2022/crowdstats/models"
)

type LoadRunner struct {
	ID         int
	Err        error
	DataDigest *models.DataDigest
	InvalidUrl string
	Loader     dataloader.Loader
	SrcUrl     string
}

func NewLoadRunner(id int, loader dataloader.Loader, srcUrl string) concurrency.Task {
	return &LoadRunner{
		ID:     id,
		SrcUrl: srcUrl,
		Loader: loader,
	}
}

// Run runs a Task and mark done using sync.WorkGroup.
func (runner *LoadRunner) Run(wg *sync.WaitGroup) {
	runner.DataDigest, runner.Err = runner.loadData(runner.Loader, runner.SrcUrl)
	wg.Done()
}

func (runner *LoadRunner) loadData(loader dataloader.Loader, srcUrl string) (*models.DataDigest, error) {
	// log.Println("task", task.ID, "started for", srcUrl)

	digest, err := loader.LoadDigest(srcUrl)
	if err != nil {
		// log.Println("task", task.ID, "errored out for", srcUrl)
		return nil, err
	}

	// log.Println("task", task.ID, "finished for", srcUrl)
	return digest, nil
}
