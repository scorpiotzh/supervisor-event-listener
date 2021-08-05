package main

import (
	"flag"
)

var (
//log = mylog.NewLogger("main", mylog.LevelDebug)
)

func main() {
	var key = flag.String("key", "", "notify key")
	flag.Parse()
	for {
		Start(*key)
	}
}
