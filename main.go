package main

import (
	"flag"
)

const (
	ResultReady = "READY\n"
	ResultOk    = "RESULT 2\nOK"
	ResultFail  = "RESULT 4\nFAIL"
)

func main() {
	var key = flag.String("key", "", "notify key")
	flag.Parse()
	for {
		Start(*key)
	}
}
