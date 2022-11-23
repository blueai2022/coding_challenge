# VirusTotal Clone Design Submitted By Ben Lu

    - Realtime UX (Ex: React JS -> gRPC or gRPC Gateway HTTP)
    - Public API (Programming access)

## Services & their APIs

    - Malware Scan Service 
        gRPC
        gRPC Gataway (HTTP)
        CLI
    
    - Scan Result Retrieval Service
        gRPC
        gRPC Gataway (HTTP)
        CLI

### Handling of median age not representing actual record:

- Design: add an extra return value for calculation function to let caller know if median value represents an actual record. Specifically for even number of records, isActual = true if middle two numbers in "sorted list" are equal. Otherwise, isActual = false. It is a byproduct of calculating median; so it is natural to return from Calc().

Implementation of MedianAge's Calc() returns three outputs:
    
    func (mdn *MedianAge) Calc() (median float64, isActual bool, err error) 

If returned isActual is false, caller Summarizer does not try to find person with median age inside "age->person name" map. It explicitly sets summary.IsMedianAgeActual = isActual. Finanlly, main() uses the summary.IsMedianAgeActual: if not actual, it explicitly outputs "median age person: < no record found >" to make distinction from empty name record (if validation was not implemented correctly).    
    
### Performance considerations:

Judging from requirement doc, I think low-latency and ability to scale is central to this project. So below are some performance-oriented implementation designs:

- medianage.go: implementation of median calculation takes advantage of ages being range-bound (between 0 to 200). 200 to handle all possible ages. General solutions such as "two heap" have relatively high time and space complexity. Below are the comparision between the implementation and popular "two heap" approach:

1. Calc() method:
    - Two Heap implementation:
        - time complexity:  5 * Olog(n) ~ O(logn)
        - space complexity: O(n)

    - Current implementation using [201]int64:
        - time complexity:  O(201) ~ O(1) 
        - space complexity: O(201) ~ O(1)

2. Add()/AddAll() method:
    - Two Heap implementation:
        - time complexity:  O(n)
        - space complexity: O(n)

    - Current implementation using [201]int64:
        - time complexity:  O(n) + O(k) ~ O(n) counting the time to populate and aggregate []DataDigest
        - space complexity: O(201) + O(k) ~ O(k) 
        - in above, k is the number of valid files, far smaller than n 

- datadigest.go & aggregate.go: implementation of DataDigest (intermediary data struct between []Person and final summary stats) takes advantge of most of the Person data are not needed for calculating the summary stats. In addition to space complexity reduction to O(403) ~ O(1), it still allows aggregation from independent and paralell data loading processes. Aggregation logic uses greedy algo "divide & conquer" to reduce time cost even more.

- validator.go: implentation avoids costly regex for validating english person names. Instead, it use an array [256]bool to check all valid name characters. 

## Testing implementation and tests performed:

1. Overall, project has 2 implementations of local test file loading logic to help tests:

- As "mock" loader which acts just like a HTTPLoader (dataloader/mock/file.go)
    FileLoader and HTTPLoader both implements the same Loader interface and provide data struct loading logic so that main code can be tested transparently   

    Used for: ./main_test.go

- As simple []byte loading logic inside handlerFunc for httptest.Server ()

    Used for: ./domain/summarizer_test.go

2. Expected results for tests are manually prepared using Excel. Ex: sort by age to find "middle" age numbers in aggregated datasets

3. Different tests performed:

- ./main_test.go: test whole flow with mock file loader and local files
    - command: "go test -v ./"

    - datasets: (source urls) 
        - ./testdata/file_src_urls.txt
        - ./testdata/file_src_urls_bad.txt
        - ./testdata/file_src_urls_dupe.txt

- ./domain/summarizer_test.go：test http data loading and summarization logic using httptest.Server
    - command: "go test -v ./domain"

    - datasets: (source urls) 
        - ./testdata/http_src_urls.txt
        - ./testdata/http_src_urls_bad.txt
        - ./testdata/http_src_urls_dupe.txt
        - scenario-based testing: dynamically created file sets inside ./testdata/generated (dir cleaned before and after each test)

- /stats/medianage_test.go: test median age logic with test data
    - command: "go test -v ./stats"
    - testing age lists are always randomly shuffled before running any tests

## How to build, test and run

Project uses Go version 1.19. Main reason for this limitation: cannot test other Go versions due to limited time on the project. The following list commands to build test and run project. They are also listed in ./Makefile:

- Run build:

    ```bash
    go build -o crowdstats main.go
    ```

- Run all tests:

    ```bash
    go test -v -cover ./...
    ```

- Run test to check output and performance of program (for test file set in zip):

    ```bash
    go test -v -cover -count=1 -run TestOriginalSet
    ```

- Run with acutal HTTP source urls file (to be supplied by tester):

    ```bash
    go run main.go <source_urls_file_path>
    ```
    
## Answers to Followup Questions

1. What assumptions did you make in your design? Why?
    - Overall reason for some assumptions made: to limit scope and keep quality high for core code modules and performance-focused code tuning

    - Person data validation: need to validate character set for English names and age needs to be an int in the correct range. Reason: seems to be reasonable and likely requirements
    
    - Unicode Person Names Validation: accepts Unicode names without character set validation. Reason: seems to be reasonable and likely requirements

    - Data file layout validation: needs to validate data having exactly 3 columns. However, column order not to be enforced. Reason: seems to be reasonable real-world requirements

    - HTTPS vs HTTP file URLs testing: tests with HTTP (not HTTPS) urls are sufficient to test HTTP data loader.
        Reason: limit scope. Or will have to install more widely accepted SSL cert to get around below error on local Mac

            x509: “Acme Co” certificate is not trusted

    - HTTP loader setting: Max of 10 HTTP redirects and 10 minute timeouts are able to handle real-world large file loading over HTTP. Reason: no time to research further on the subject

    - Maximum of 10 worker threads is okay for local testing. Reason: think this will be tested on other engieer's local dev Mac. Also, scaling is discussed in answers to follow-up questions #2 and #3 
    
    - Requirement Interpretations:
        - No need to check for file content that is duplicated inside multiple data files. Reason: limit scope & not in requirement doc
        
        - No need to fail a summary run even if many files are invalid. Reason: not in requirement doc

    - Although testing for median age can be further enhanced with random test data generation, but I assume it is not necessary for this take-home project. Reason: limit scope and implementation is relatively trivial. It is also tested sufficiently with tests written in ./main_test.go and ./domain/summarizer_test.go

    Random test data generation for median age: use 1 or 2 random generated middle number(s) as expected answer(s); then pad equal (random) count of random numbers that are <= and >= the chosen middle number(s). (Ex: random single middle number: 35. Pad 25 random numbers <= 35 and 25 random numbers >=35 to form final list)

2. How would you change your program if it had to process many files where each file was over 10M records?

    - Basic idea: Base "unit of work" (that can be run in parallel and its result later aggregated with others'): 
        - Code location: func Digest(io.Reader) inside ./dataloader.digester.go
        
        - If input io.Reader represents a subset of records (rather than 10M), we can allocate a separate worker pool to process it in parallel as if we do for many different files.

    - Option 1 - to implement on a single server:
        - Reduce the size (Ex: size of 5) of current worker pool for HTTPLoader (call it "HTTPLoader Pool"). Its count represents how many files can be processed at one time. It runs HTTPLoader data load logic before calling Digest()

        - Modify HTTPLoader to read the datastrem of a single file and collect a subset of (Ex: 50K) records; create a new StringReader using the subset data and pass it to be run with Digest() (running on separate thread pool)

        - The worked pool to run Digest() needs to much larger (Ex: size of 25 or as many as the server can afford without performance penalty) so that it can serve many tasks from any of HTTPLoaders running in parallel. Called "Digest pool or threads" in this section.

        - This implementation maximizes the processing power of a single server. For example, it will use all available "Digest threads" for a large file rather than only 1 thread per file: 
            - subsets per 10M record file: 10M / 50K = 200
            - workload per Digest thread: 200 / 25 = 8

    - Option 2 - to implement Digest() logic as client streaming gRPC service on a Kubernetts cluster: 
        - The whole cluster's processing power, rather than a single server, is available to scale up performance. Aggregation on behalf of a single file is simple, since "client streaming" gRPC server only returns one response message for numerous request messages on behalf of a single large file. It's important to "correlate" responses to a file in this way since we want to fail the entire file if any chunk is invalid.

        - Other modifications: Digest() gRPC call takes a list of strings rather than Reader, 

    - Option 2 is more preferrable with slightly more complexity; gRPC over HTTP2 uses binary format and message compression so input parameter size is also okay for this option

3. How would you change your program if it had to process data from more than 20K URLs?

    - Assumption: to make this solution somewhat differentiated than answer to Follow-up Question #2, I assume that we are mainly deal with a large number of files to minimize our implementation changes. If both file sizes and number of files are large: need to use a combination of both designs

    - If file size is not much of a concern; it is more efficient to keep HTTPLoader loading logic and Digest() logic together as a single unit since it is quite memory efficient: it increments the DataDigest through reading the data streams.

    - Option 1 - to implement on a single server:
        - We need tune thread pool to fully utilize production server processing capacity

        - Other than that, we don't need much changes to the current implemention: keeping HTTPLoader loading logic and Digest() logic together as a single unit

    - Option 2: implement HTTPLoader and Digest together as "backend" unary gRPC server:
        - Serving the intial client call, we can keep the Summarizer.Run() logic as is. Only to modify logic inside loaderrunner.go to make the backend unary gRPC server call with a file URL.

        - Increase original thread pool size for local Summarizer.Run() to be larger (Ex: 2 x # of cores) and there is not much CPU or memory load locally. Each threaded execution waits for unary gRPC call to complete.

        - HTTPLoader and Digest stays together as "backend" gRPC server to load and digest one file at a time. The Kubernetts cluster's processing power, rather than a single server, is available to scale up performance. A small number of clustered pods can easily handle more than 100 concurrent requests at a time.

        - Locally, Summarizer.Run() will run as is: it waits for all 20k files to be processed before aggregating results and summarizing them

    - Option 3: implement HTTPLoader and Digest together as "backend" client streaming gRPC server:
        - The main difference from Option 2 is that message exchanges are less chatty

        - To use this option, loadrunner needs to be modified to take a small batch of urls to send as a sequence of messages  with the same stream to gRPC server
    
    - Choose option 3 for a performant distributed system
