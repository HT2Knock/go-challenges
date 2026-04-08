package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type ToDoService struct {
	requestGroup singleflight.Group
	cache        sync.Map
}

func (p *ToDoService) fetchToDoApi(task string) (string, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/" + task)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func (p *ToDoService) GetToDo(task string) (string, error) {
	if data, ok := p.cache.Load(task); ok {
		return data.(string), nil
	}

	v, err, _ := p.requestGroup.Do(task, func() (any, error) {
		task, err := p.fetchToDoApi(task)
		if err == nil {
			p.cache.Store(task, task)
		}

		return task, nil
	})
	if err != nil {
		return "", err
	}

	return v.(string), err
}

func main() {
	service := &ToDoService{}

	for i := 0; i < 10; i++ {
		go func(i int) {
			task, err := service.GetToDo("1")

			if err == nil {
				fmt.Printf("Goroutine %d got task data: %s\n", i, task)
			} else {
				fmt.Printf("Goroutine %d encountered an error: %s\n", i, err)
			}
		}(i)
	}

	time.Sleep(5 * time.Second)
}
