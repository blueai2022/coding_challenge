package domain

import (
	"log"
	"time"

	"github.com/blucv2022/crowdstats/concurrency"
	"github.com/blucv2022/crowdstats/dataloader"
	"github.com/blucv2022/crowdstats/models"
	"github.com/blucv2022/crowdstats/stats"
)

const (
	maxThreads       = 10
	avgJobsPerThread = 2
)

var (
	valExists = struct{}{}
)

type Summarizer struct {
	loader dataloader.Loader
}

func NewSummarizer(loader dataloader.Loader) *Summarizer {
	return &Summarizer{
		loader: loader,
	}
}

func (smr *Summarizer) Run(srcUrls []string) (*models.Summary, *models.RunStats) {
	origCount := len(srcUrls)
	srcUrls, dupeUrls := dedupeUrls(srcUrls)
	dupeCount := origCount - len(srcUrls)

	numJobs := len(srcUrls)
	numThreads := min(maxThreads, int(numJobs+1)/avgJobsPerThread)

	var tasks []concurrency.Task

	for jobId := 1; jobId <= numJobs; jobId++ {
		tasks = append(tasks, NewLoadRunner(jobId, smr.loader, srcUrls[jobId-1]))
	}

	pool := concurrency.NewPool(tasks, numThreads)

	digests := []*models.DataDigest{}
	invalidUrls := []string{}

	//reading starts
	startTime := time.Now()
	pool.Run()

	for _, task := range tasks {
		loadRunner, ok := task.(*LoadRunner)
		if !ok {
			log.Fatalf("unexpected Task from pool:%v", task)
		}

		if loadRunner.Err != nil {
			if _, ok := loadRunner.Err.(*dataloader.LoadError); ok {
				invalidUrls = append(invalidUrls, loadRunner.SrcUrl)
			} else {
				log.Fatal("cannot load data d/t unexpected error: ", loadRunner.Err)
			}
		} else {
			if loadRunner.DataDigest == nil {
				log.Fatal("data digest not loaded")
			}

			if loadRunner.DataDigest.TotalAgeCounts > 0 {
				digests = append(digests, loadRunner.DataDigest)
			}
		}
	}

	//reading ends
	readTimeCost := time.Since(startTime)

	summary := smr.fromDigests(digests)

	runStats := &models.RunStats{
		ReadTimeCost:      readTimeCost,
		InvalidUrls:       invalidUrls,
		DuplicateUrls:     dupeUrls,
		DuplicateUrlCount: dupeCount,
		NumThreads:        numThreads,
	}

	return summary, runStats
}

func (smr *Summarizer) fromDigests(digests []*models.DataDigest) *models.Summary {
	if len(digests) == 0 {
		return &models.Summary{
			IsNA: true,
		}
	}

	aggrDigest := smr.aggregateK(digests)

	medianCalculator := stats.NewMedianAge()
	medianCalculator.AddAll(aggrDigest.AgeCounts, aggrDigest.TotalAgeCounts)

	medianAge, isActual, err := medianCalculator.Calc()
	if err != nil {
		log.Fatal("unexpected error from calculating median:", err)
	}

	var medianAgePerson string
	if isActual {
		medianAgePerson = aggrDigest.AgePersonName[int(medianAge)]
	}

	avgAge, err := stats.AverageAge(aggrDigest.AgeCounts)
	if err != nil {
		log.Fatal("cannot calculate average age:", err)
	}

	return &models.Summary{
		IsNA:              false,
		AverageAge:        float32(avgAge),
		MedianAge:         float32(medianAge),
		MedianAgePerson:   medianAgePerson,
		IsMedianAgeActual: isActual,
	}
}

func dedupeUrls(srcUrls []string) ([]string, []string) {
	validUrls := []string{}
	dupeUrls := []string{}

	urlSet := make(map[string]struct{})

	for _, url := range srcUrls {
		if _, ok := urlSet[url]; !ok {
			validUrls = append(validUrls, url)
			urlSet[url] = valExists
		} else {
			dupeUrls = append(dupeUrls, url)
		}
	}

	return validUrls, dupeUrls
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
