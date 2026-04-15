package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
)

type LogEntryV1 struct {
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

type LogEntryV2 struct {
	CreatedAt time.Time `parquet:"created_at,timestamp"`
	Actor     *struct {
		Login      string `parquet:"login,string"`
		GravatarID string `parquet:"gravatar_id,string"`
		URL        string `parquet:"url,string"`
		AvatarURL  string `parquet:"avatar_url,string"`
		ID         int    `parquet:"id,int64"`
	} `parquet:"actor,optional"`
	Repo *struct {
		Name string `parquet:"name,string"`
		URL  string `parquet:"url,string"`
		ID   int    `parquet:"id,int64"`
	} `parquet:"repo,optional"`
	Payload *struct {
		Ref          string `parquet:"ref,string"`
		RefType      string `parquet:"ref_type,string"`
		MasterBranch string `parquet:"master_branch,string"`
		Description  string `parquet:"description,string"`
		PusherType   string `parquet:"pusher_type,string"`
	} `parquet:"payload,optional"`
	ID     string `parquet:"id,string"`
	Type   string `parquet:"type,dict,string"`
	Public bool   `parquet:"public,boolean"`
}

func convert(path string) error {
	ext := filepath.Ext(path)
	outPath := path[:len(path)-len(ext)] + ".parquet"

	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening json: %w", err)
	}
	defer jsonFile.Close()

	pf, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating parquet: %w", err)
	}
	defer pf.Close()

	writer := parquet.NewGenericWriter[LogEntryV1](pf, parquet.Compression(&parquet.Zstd))

	decoder := json.NewDecoder(jsonFile)
	if _, err = decoder.Token(); err != nil {
		return fmt.Errorf("expected array start: %w", err)
	}

	const batchSize = 1000
	batch := make([]LogEntryV1, 0, batchSize)

	for decoder.More() {
		var l LogEntryV1
		if err := decoder.Decode(&l); err != nil {
			return fmt.Errorf("decode error: %w", err)
		}

		batch = append(batch, l)

		if len(batch) >= batchSize {
			if _, err := writer.Write(batch); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
			batch = batch[:0]
		}

	}

	if len(batch) > 0 {
		if _, err := writer.Write(batch); err != nil {
			return fmt.Errorf("final write batch error: %w", err)
		}
	}

	return writer.Close()
}

func convertV2(path string) error {
	ext := filepath.Ext(path)
	outPath := path[:len(path)-len(ext)] + ".parquet"

	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening json: %w", err)
	}
	defer jsonFile.Close()

	// Buffer reader and writer to read a larger chunk
	reader := bufio.NewReaderSize(jsonFile, 32*1024)

	pf, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating parquet: %w", err)
	}
	defer pf.Close()

	writer := parquet.NewGenericWriter[LogEntryV2](pf, parquet.Compression(&parquet.Zstd))

	decoder := json.NewDecoder(reader)
	if _, err = decoder.Token(); err != nil {
		return fmt.Errorf("expected array start: %w", err)
	}

	const batchSize = 1000
	batch := make([]LogEntryV2, 0, batchSize)

	var l LogEntryV2
	for decoder.More() {
		l = LogEntryV2{}

		if err := decoder.Decode(&l); err != nil {
			return fmt.Errorf("decode error: %w", err)
		}

		batch = append(batch, l)

		if len(batch) >= batchSize {
			if _, err := writer.Write(batch); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
			batch = batch[:0]
		}

	}

	if len(batch) > 0 {
		if _, err := writer.Write(batch); err != nil {
			return fmt.Errorf("final write batch error: %w", err)
		}
	}

	return writer.Close()
}

func main() {
	if err := convertV2("./data/large-file.json"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully converted json to parquet")
}
