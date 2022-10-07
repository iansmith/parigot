package ex1

import (
	"fmt"
	"github.com/iansmith/parigot/lib/base/go/log"
)

type composite struct {
	name string
	age  int
}

func sayHi(c *composite) string {
	return "hi " + c.name + ", when do you turn " + fmt.Sprint(c.age+1) + "?\n"
}

func Driver() {
	print("0\n")
	c := &composite{"billy", 13}
	log.Dev.Debug(sayHi(c))
}
