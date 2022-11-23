package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blucv2022/crowdstats/dataloader"
	"github.com/blucv2022/crowdstats/domain"
	"github.com/blucv2022/crowdstats/models"
)

var (
	exists = struct{}{}
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("source url file not provided - Usage \"go run main.go <source_urls_file_path>\"")
	}

	srcUrlFile := os.Args[1]
	loader := dataloader.NewHttpLoader()

	batch(srcUrlFile, loader)
}

func batch(srcUrlFile string, loader dataloader.Loader) (summary *models.Summary, runStats *models.RunStats) {
	srcUrls := getDataSrcUrls(srcUrlFile)

	summarizer := domain.NewSummarizer(loader)
	summary, runStats = summarizer.Run(srcUrls)
	outputResult(summary, runStats)

	return
}

func getDataSrcUrls(filePath string) []string {
	// open source urls file
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("cannot open source url file:", err)
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
		log.Fatal("cannot read source url lines:", err)
	}

	return res
}

func outputResult(summary *models.Summary, runStats *models.RunStats) {
	if summary.IsNA {
		fmt.Println("dataset summary not available: data not found")
		return
	}

	fmt.Println("dataset summary:")
	fmt.Printf("  average age: %.2f\n", summary.AverageAge)
	fmt.Printf("  median age: %.2f\n", summary.MedianAge)

	if summary.IsMedianAgeActual {
		fmt.Printf("  median age person: %s\n", summary.MedianAgePerson)
	} else {
		fmt.Println("  median age person: <no record found>")
	}

	fmt.Println()
	fmt.Printf("dataset read threads: %d\n", runStats.NumThreads)
	fmt.Printf("dataset read time cost: %v\n", runStats.ReadTimeCost)

	fmt.Println()
	fmt.Printf("rejected data files count: %d\n", len(runStats.InvalidUrls))
	for _, url := range runStats.InvalidUrls {
		fmt.Printf("  - %s\n", url)
	}

	fmt.Println()
	fmt.Printf("duplicate data files count: %d\n", runStats.DuplicateUrlCount)
	for _, url := range runStats.DuplicateUrls {
		fmt.Printf("  - %s\n", url)
	}
}
