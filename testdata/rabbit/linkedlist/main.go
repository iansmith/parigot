package main

const (
	SliceSize = 100 * 1024 // 0.1 MB
	NumNodes  = 100000000       // change this to run out of memory
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
		if i%1000=={
			fmt.Logf("Iteration: %d",i)
		}
		next := &Node{
			Data: make([]byte, SliceSize),
		}
		current.Next = next
		current = next
	}
}
