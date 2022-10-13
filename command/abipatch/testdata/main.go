package main

type big struct {
	a int64
	b int64
	c int64
}

type foo struct {
	b    big
	t    int64
	frak byte
	frik int64
}

//log.Dev.Debug("hello, logger")
//ex1.Driver()
//log.Dev.Debug("goodbye, logger")

func main() {
	print(zap([]byte("foobie")))
}

//export zap
func zap(b []byte) *byte {
	return &b[0]
}
