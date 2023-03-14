package pbmodel

import (
	"log"

	"github.com/dominikbraun/graph"
)

var Pb3Dep = graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

func AddImportEdge(from, to string) {
	Pb3Dep.AddVertex(from)
	Pb3Dep.AddVertex(to)
	Pb3Dep.AddEdge(from, to)
	log.Printf("added edge: %s -> %s", from, to)
}
