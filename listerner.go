package main

import (
	"bufio"
	"errors"
	"log"
	"os"
)

var (
	ErrPayloadLength = errors.New("Header中len长度与实际读取长度不一致")
)

func Start(key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()
	listen(key)
}

// 监听事件, 从标准输入获取事件内容
func listen(key string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Println("READY ...")
		header, err := readHeader(reader)
		if err != nil {
			log.Println("readHeader err:", err.Error())
			continue
		}
		payload, err := readPayload(reader, header.Len)
		if err != nil {
			log.Println("readPayload err:", err.Error())
			continue
		}
		msg := Message{Header: header, Payload: payload}
		switch header.EventName {
		case "PROCESS_STATE_EXITED", "PROCESS_STATE_BACKOFF", "PROCESS_STATE_STOPPED", "PROCESS_STATE_FATAL":
			SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		case "PROCESS_STATE_STARTING", "PROCESS_STATE_UNKNOWN", "PROCESS_STATE_STOPPING":
		case "PROCESS_STATE_RUNNING":
			SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		default:
			SendLarkTextNotify(key, "程序状态变化事件通知", msg.String())
		}
		log.Println("SUCCESS ...")
	}
}

// 读取header
func readHeader(reader *bufio.Reader) (*Header, error) {
	// 读取Header
	data, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// 解析Header
	header, err := ParseHeader(data)
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
