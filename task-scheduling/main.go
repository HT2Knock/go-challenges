package main

import (
	"fmt"
	"math/rand"
	"slices"
)

type Task struct {
	id           string
	dependencies []string
	durationMs   int
}

func generateTask(n int) []Task {
	tasks := make([]Task, n)

	for i := range tasks {
		id := fmt.Sprintf("task-%d", i)

		var deps []string

		if i > 0 {
			numDeps := rand.Intn(3)

			for range numDeps {
				randomPrevIndex := rand.Intn(i)
				depID := fmt.Sprintf("task-%d", randomPrevIndex)

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
	g := NewGraph()

	for _, t := range tasks {
		fmt.Printf("ID: %s, Dur: %d, Deps: %v \n", t.id, t.durationMs, t.dependencies)
		g.vertices[t.id] = append(g.vertices[t.id], t.dependencies...)
	}

	fmt.Println(g.vertices)
}
