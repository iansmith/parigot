package lib

import (
	"fmt"
)

type TypeAnnotation struct {
	Number int
}

func (t *TypeAnnotation) String() string {
	return fmt.Sprintf("(;%d;)", t.Number)
}
