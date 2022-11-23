build:
	go build -o crowdstats main.go

test:
	go test -v -cover ./...

test_main:
	go test -v -cover -count=1 -run TestOriginalSet

run:
	go run main.go $(source_urls_file)

.PHONY: build, test, test_main, run
