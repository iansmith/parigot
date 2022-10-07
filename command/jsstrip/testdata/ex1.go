package main

import "fmt"

type composite struct {
	name string
	age  int
}

func sayHi(c *composite) string {
	return "hi " + c.name + ", when do you turn " + fmt.Sprint(c.age+1) + "\n"
}

func driver() {
	print("0\n")
	c := &composite{"billy", 13}
	print("1=n")
	print(sayHi(c))
}
