package ex1

import (
	"fmt"
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
	_ = &composite{"billy", 13}
	//log.Dev.Debug(sayHi(c))
}
