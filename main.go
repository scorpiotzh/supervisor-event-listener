package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ResultReady = "READY\n"
	ResultOk    = "RESULT 2\nOK"
	ResultFail  = "RESULT 4\nFAIL"
)

func main() {
	var key = flag.String("key", "", "notify key")
	flag.Parse()
	fmt.Println("key:", *key)

	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	for {
		// 发送后等待接收event
		_, _ = stdout.WriteString(ResultReady)
		_ = stdout.Flush()
		// 接收header
		line, _, _ := stdin.ReadLine()
		_, _ = stderr.WriteString("read" + string(line))
		_ = stderr.Flush()

		header, payloadSize := parseHeader(line)

		// 接收payload
		payload := make([]byte, payloadSize)
		_, _ = stdin.Read(payload)
		_, _ = stderr.WriteString("read : " + string(payload))
		_ = stderr.Flush()

		result := alarm(header, payload)

		if result { // 发送处理结果
			_, _ = stdout.WriteString(ResultOk)
		} else {
			_, _ = stdout.WriteString(ResultFail)
		}
		_ = stdout.Flush()
	}
}

func parseHeader(data []byte) (header map[string]string,
	payloadSize int) {
	pairs := strings.Split(string(data), " ")
	header = make(map[string]string, len(pairs))

	for _, pair := range pairs {
		token := strings.Split(pair, ":")
		header[token[0]] = token[1]
	}

	payloadSize, _ = strconv.Atoi(header["len"])
	return header, payloadSize
}

// 这里设置报警即可
func alarm(header map[string]string, payload []byte) bool {
	// send mail
	fmt.Println(header, string(payload))
	return true
}
