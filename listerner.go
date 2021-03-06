package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var (
	ErrPayloadLength = errors.New("header 中 len 长度与实际读取长度不一致")
	stdin            = bufio.NewReader(os.Stdin)
	stdout           = bufio.NewWriter(os.Stdout)
	stderr           = bufio.NewWriter(os.Stderr)
)

func Start(key string) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = stderr.WriteString(fmt.Sprintf("panic: %v", err))
			_ = stderr.Flush()
		}
	}()
	listen(key)
}

// 监听事件, 从标准输入获取事件内容
func listen(key string) {
	for {
		// 发送后等待接收 event
		ready()
		// 接收 header
		header, err := readHeader(stdin)
		if err != nil {
			failure("readHeader", err)
			continue
		}
		// 接收 payload
		payload, err := readPayload(stdin, header.Len)
		if err != nil {
			failure("readPayload", err)
			continue
		}
		msg := Message{Header: header, Payload: payload}
		switch header.EventName {
		case "PROCESS_STATE_EXITED", "PROCESS_STATE_BACKOFF", "PROCESS_STATE_STOPPED", "PROCESS_STATE_FATAL":
			_ = SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		case "PROCESS_STATE_STARTING", "PROCESS_STATE_UNKNOWN", "PROCESS_STATE_STOPPING":
		case "PROCESS_STATE_RUNNING":
			_ = SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		default:
			_ = SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		}
		success()
	}
}

// 读取header
func readHeader(reader *bufio.Reader) (*Header, error) {
	// 读取Header
	data, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	// 解析Header
	header, err := ParseHeader(string(data))
	if err != nil {
		return nil, err
	}

	return header, nil
}

// 读取payload
func readPayload(reader *bufio.Reader, payloadLen int) (*Payload, error) {
	// 读取payload
	buf := make([]byte, payloadLen)
	length, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if payloadLen != length {
		return nil, ErrPayloadLength
	}
	// 解析payload
	payload, err := ParsePayload(string(buf))
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func ready() {
	_, _ = stdout.WriteString(ResultReady)
	_ = stdout.Flush()
}

func success() {
	_, _ = stdout.WriteString(ResultOk)
	_ = stdout.Flush()
}

func failure(funcName string, err error) {
	_, _ = stderr.WriteString(funcName + ": \n" + err.Error())
	_ = stderr.Flush()
	_, _ = stdout.WriteString(ResultFail)
	_ = stdout.Flush()
}
