package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
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

func main() {
    // Define flags
    esHost := flag.String("es_host", "http://localhost:9200", "Elasticsearch host URL")
    username := flag.String("username", "", "Elasticsearch username")
    password := flag.String("password", "", "Elasticsearch password")
    indexName := flag.String("indexname", "my-index", "Elasticsearch index name")
    tool := flag.String("tool", "tool", "Tool name")
    testFlag := flag.Bool("test", false, "Test Elasticsearch connection")

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

    // For normal execution, read from stdin and insert each line as a document
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        doc := Document{
            Date:     time.Now().Format(time.RFC3339),
            Output:   line,
            ToolName: *tool,
        }

        err := insertDocument(es, *indexName, doc)
        if err != nil {
            log.Printf("Failed to insert document: %s", err)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading from stdin: %s", err)
    }
}

