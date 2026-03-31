package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"slices"
)

type Task struct {
	id           string
	dependencies []string
	durationMs   int
}

func topSort(tasks []Task) []string {
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)
	queue := list.New()
	sortedList := []string{}

	for _, task := range tasks {
		inDegree[task.id] = len(task.dependencies)

		for _, dep := range task.dependencies {
			adjList[dep] = append(adjList[dep], task.id)
		}
	}

	for _, task := range tasks {
		if inDegree[task.id] == 0 {
			queue.PushBack(task.id)
		}
	}

	for queue.Len() > 0 {
		e := queue.Front()
		u := e.Value.(string)
		queue.Remove(e)

		sortedList = append(sortedList, u)

		for _, neighbor := range adjList[u] {
			inDegree[neighbor] -= 1
			if inDegree[neighbor] == 0 {
				queue.PushBack(neighbor)
			}
		}
	}

	return sortedList
}

func generateTask(n int) []Task {
	tasks := make([]Task, n)

	for i := range tasks {
		id := fmt.Sprintf("task-%d", i)

		var deps []string

		if i > 0 {
			for range rand.Intn(3) {
				depID := fmt.Sprintf("task-%d", rand.Intn(i))

				if !slices.Contains(deps, depID) {
					deps = append(deps, depID)
				}
			}
		}

		tasks[i] = Task{
			id:           id,
			durationMs:   rand.Intn(5000),
			dependencies: deps,
		}
	}

	return tasks
}

func main() {
	tasks := generateTask(10)

	for _, t := range tasks {
		fmt.Printf("ID: %s, Dur: %d, Deps: %v \n", t.id, t.durationMs, t.dependencies)
	}

	sortedList := topSort(tasks)

	fmt.Println(sortedList)
}
