package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/elastic/go-elasticsearch/v8"
)

type Document struct {
    Date     string `json:"date"`
    Output   string `json:"output"`
    ToolName string `json:"toolname"`
}

func insertDocument(es *elasticsearch.Client, index string, doc Document) error {
    docJSON, err := json.Marshal(doc)
    if err != nil {
        return fmt.Errorf("failed to marshal document: %w", err)
    }

    res, err := es.Index(index, strings.NewReader(string(docJSON)))
    if err != nil {
        return fmt.Errorf("failed to insert document: %w", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("failed to insert document: %s", res.String())
    }

    log.Printf("Document inserted into index %s", index)
    return nil
}

func worker(es *elasticsearch.Client, index string, toolName string, lines <-chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    for line := range lines {
        doc := Document{
            Date:     time.Now().Format(time.RFC3339),
            Output:   line,
            ToolName: toolName,
        }
        if err := insertDocument(es, index, doc); err != nil {
            log.Printf("Failed to insert document: %s", err)
        }
    }
}

func main() {
    // Define flags
    esHost := flag.String("es_host", "http://localhost:9200", "Elasticsearch host URL")
    username := flag.String("username", "", "Elasticsearch username")
    password := flag.String("password", "", "Elasticsearch password")
    indexName := flag.String("indexname", "my-index", "Elasticsearch index name")
    tool := flag.String("tool", "tool", "Tool name")
    testFlag := flag.Bool("test", false, "Test Elasticsearch connection")
    numWorkers := flag.Int("workers", 5, "Number of concurrent workers")

    flag.Parse()

    // Configure Elasticsearch client
    cfg := elasticsearch.Config{
        Addresses: []string{*esHost},
        Username:  *username,
        Password:  *password,
    }
    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
        log.Fatalf("Error creating Elasticsearch client: %s", err)
    }

    // Test connection if --test flag is set
    if *testFlag {
        log.Println("Testing Elasticsearch connection by inserting sample data...")

        // Create sample document
        sampleDoc := Document{
            Date:     time.Now().Format(time.RFC3339),
            Output:   "Sample output for testing",
            ToolName: "sample-tool",
        }

        // Attempt to insert into a sample index
        err := insertDocument(es, "sample-index", sampleDoc)
        if err != nil {
            log.Fatalf("Test failed: %s", err)
        }
        log.Println("Test succeeded: Sample document inserted into 'sample-index'")
        return
    }

    // Create a channel to send lines to workers
    lines := make(chan string)
    var wg sync.WaitGroup

    // Start the worker pool
    for i := 0; i < *numWorkers; i++ {
        wg.Add(1)
        go worker(es, *indexName, *tool, lines, &wg)
    }

    // Read from stdin and send lines to the workers immediately
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        lines <- scanner.Text()
    }
    close(lines) // Close channel to signal workers to finish

    // Wait for all workers to complete
    wg.Wait()

    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading from stdin: %s", err)
    }
}
