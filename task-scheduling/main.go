package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type Task struct {
	id           string
	dependencies []string
	durationMs   int
}

// minTotalTime using top sort Kahn algorithms
func minTotalTime(tasks []Task) (int, error) {
	inDegree := make(map[string]int)
	adjList := make(map[string][]Task)

	for _, task := range tasks {
		inDegree[task.id] = len(task.dependencies)

		for _, dep := range task.dependencies {
			adjList[dep] = append(adjList[dep], task)
		}
	}

	queue := list.New()
	earliestFinish := make(map[string]int)

	for _, task := range tasks {
		if inDegree[task.id] == 0 {
			queue.PushBack(task)
			earliestFinish[task.id] = task.durationMs
		}
	}

	sortedList := []string{}

	for queue.Len() > 0 {
		e := queue.Front()
		u := e.Value.(Task)
		queue.Remove(e)

		sortedList = append(sortedList, u.id)
		uFinishTime := earliestFinish[u.id]

		for _, neighbor := range adjList[u.id] {
			newFinishTime := uFinishTime + neighbor.durationMs
			if newFinishTime > earliestFinish[neighbor.id] {
				earliestFinish[neighbor.id] = newFinishTime
			}

			inDegree[neighbor.id]--
			if inDegree[neighbor.id] == 0 {
				queue.PushBack(neighbor)
			}
		}
	}

	if len(sortedList) != len(tasks) {
		return 0, errors.New("cycle detected")
	}

	maxTime := 0
	for _, t := range earliestFinish {
		if t > maxTime {
			maxTime = t
		}
	}

	return maxTime, nil
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
		fmt.Printf("ID: %s, Dur: %d ms, Deps: %v \n", t.id, t.durationMs, t.dependencies)
	}

	minTotal, err := minTotalTime(tasks)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Minimum total time for task scheduling = %d ms", minTotal)
}
