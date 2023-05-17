package main

const (
	SliceSize = 4096 * 1024 // 4096KB
	NumNodes  = 10000       // change this to run out of memory
)

type Node struct {
	Data []byte
	Next *Node
}

func main() {
	// Create the initial node
	head := &Node{
		Data: make([]byte, SliceSize),
	}

	// Allocate additional nodes
	current := head
	for i := 0; i < NumNodes; i++ {
		next := &Node{
			Data: make([]byte, SliceSize),
		}
		current.Next = next
		current = next
	}
}
