package main

import (
	"github.com/yourbasic/graph"
)

const (
	White = iota
	Gray
	Black
)

type DFSData struct {
	Time     int
	Color    []int
	Prev     []int
	Discover []int
	Finish   []int
}

func DFS(g graph.Iterator) DFSData {
	n := g.Order() // Order returns the number of vertices.
	d := DFSData{
		Time:     0,
		Color:    make([]int, n),
		Prev:     make([]int, n),
		Discover: make([]int, n),
		Finish:   make([]int, n),
	}
	for v := 0; v < n; v++ {
		d.Color[v] = White
		d.Prev[v] = -1
	}
	for v := 0; v < n; v++ {
		if d.Color[v] == White {
			d.dfsVisit(g, v)
		}
	}
	return d
}

func (d *DFSData) dfsVisit(g graph.Iterator, v int) {
	d.Color[v] = Gray
	d.Time++
	d.Discover[v] = d.Time
	// Visit calls a function for each neighbor w of v,
	// with c equal to the cost of the edge (v, w).
	// The iteration is aborted if the function returns true.
	g.Visit(v, func(w int, c int64) (skip bool) {
		if d.Color[w] == White {
			d.Prev[w] = v
			d.dfsVisit(g, w)
		}
		return
	})
	d.Color[v] = Black
	d.Time++
	d.Finish[v] = d.Time
}
