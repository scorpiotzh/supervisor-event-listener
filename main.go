package main

import (
	"flag"
	"github.com/scorpiotzh/mylog"
)

var (
	log = mylog.NewLogger("main", mylog.LevelDebug)
)

func main() {
	var key = flag.String("key", "", "notify key")
	flag.Parse()
	for {
		Start(*key)
	}
}
