package pbmodel

import (
	"log"

	"github.com/dominikbraun/graph"
)

var pb3Import = graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

func AddImportEdge(from, to string) {
	pb3Import.AddVertex(from)
	pb3Import.AddVertex(to)
	pb3Import.AddEdge(from, to)
	log.Printf("added edge: %s -> %s", from, to)
}
