package main

type Graph struct {
	vertices map[string][]string
}

func (g *Graph) AddEdge(v, w string) {
	g.vertices[v] = append(g.vertices[v], w)
}

func NewGraph() *Graph {
	return &Graph{vertices: make(map[string][]string)}
}

func (g *Graph) bfs(start string) {
}
