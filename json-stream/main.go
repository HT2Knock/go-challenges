package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
)

type LogEntry struct {
	ID    string `parquet:"id,string"`
	Type  string `parquet:"type,dict,string"`
	Actor *struct {
		ID         int    `parquet:"id,int64"`
		Login      string `parquet:"login,string"`
		GravatarID string `parquet:"gravatar_id,string"`
		URL        string `parquet:"url,string"`
		AvatarURL  string `parquet:"avatar_url,string"`
	} `parquet:"actor,optional"`
	Repo *struct {
		ID   int    `parquet:"id,int64"`
		Name string `parquet:"name,string"`
		URL  string `parquet:"url,string"`
	} `parquet:"repo,optional"`
	Payload *struct {
		Ref          string `parquet:"ref,string"`
		RefType      string `parquet:"ref_type,string"`
		MasterBranch string `parquet:"master_branch,string"`
		Description  string `parquet:"description,string"`
		PusherType   string `parquet:"pusher_type,string"`
	} `parquet:"payload,optional"`
	Public    bool      `parquet:"public,boolean"`
	CreatedAt time.Time `parquet:"created_at,timestamp"`
}

func JsonToParquet(path string) error {
	pf, err := os.Create(path[:len(path)-len(filepath.Ext(path))] + ".parquet")
	if err != nil {
		return fmt.Errorf("failed to create parquet file with err: %w", err)
	}
	defer pf.Close()

	writer := parquet.NewGenericWriter[LogEntry](pf, parquet.Compression(&parquet.Zstd))

	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file with err: %w", err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	const batchSize = 1000
	logs := make([]LogEntry, 0, batchSize)

	for decoder.More() {
		var l LogEntry
		if err := decoder.Decode(&l); err != nil {
			return fmt.Errorf("failed to decode log with err: %w", err)
		}

		logs = append(logs, l)

		if len(logs) > batchSize {
			if _, err := writer.Write(logs); err != nil {
				return err
			}

			logs = logs[:0]
		}

	}

	if len(logs) > 0 {
		if _, err := writer.Write(logs); err != nil {
			return err
		}
	}

	_, err = decoder.Token()
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to write buffer to parquet with err: %w", err)
	}

	return nil
}

func main() {
	if err := JsonToParquet("./data/large-file.json"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully converted json to parquet")
}
