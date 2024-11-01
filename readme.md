# esdata

`esdata` is a CLI tool written in Go for inserting data into an Elasticsearch index. It reads input from standard input (`stdin`), formats each line into JSON, and stores it in the specified Elasticsearch index. This tool supports multi-threaded ingestion for efficient processing and allows for a custom tool name.

## Features

- **Multi-threaded Ingestion**: Inserts data concurrently into Elasticsearch.
- **Customizable Tool Name**: Optionally specify a tool name (default: `"tool"`).
- **Connection Testing**: Use the `--test` flag to verify the Elasticsearch connection by inserting sample data.
- **Cross-Platform Support**: Binaries available for multiple Linux architectures (amd64, arm64, 386).


### Docker (Recommended)

To build and run `esdata` using Docker:

```bash
docker build -t esdata .
```


Help menu - tool supports following options.

```
Usage of ./esdata-amd64:
  -es_host string
        Elasticsearch host URL (default "http://localhost:9200")
  -indexname string
        Elasticsearch index name (default "my-index")
  -password string
        Elasticsearch password
  -test
        Test Elasticsearch connection
  -tool string
        Tool name (default "tool")
  -username string
        Elasticsearch username
```

## Exmaple command

```
echo "Sample log entry" | ./esdata --es_host "http://localhost:9200" --username "admin" --password "admin" --indexname "logs" --tool "myTool"
```

## Esdata Example

<img width="1159" alt="image" src="https://github.com/user-attachments/assets/ccd32ba0-9168-4d17-b49d-02c2ae4cc79c">
