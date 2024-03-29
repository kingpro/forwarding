package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

var (
	listenAddr        = flag.String("listenAddr", ":8601", "listen address")
	handshakePassword = flag.String("handshakePassword", "", "encryption password. empty means no encryption")
)

type ConnInfo struct {
	conn    net.Conn // 客户端tcp连接
	chid    string   // 客户端的 channel id
	timeout uint32   // 客户端连接的心跳超时
	recvers []string // 接收者的chids
}

type SendQEle struct {
	chid string        // 发送者的 channel id
	data *bytes.Buffer // 发送的数据
}

var cmdQ = make(chan string)           // gorouting 控制命令
var sendQ = make(chan SendQEle)        // 要发送的数据发给 Dispatcher
var createConnQ = make(chan *ConnInfo) // 客户端连接创建
var removeConnQ = make(chan string)    // 客户端连接关闭

func main() {
	flag.Parse()
	log.Println("-- Main started.", *listenAddr, *handshakePassword)

	go Dispatcher()
	go Accepter(*listenAddr, *handshakePassword)

	select {
	case cmd := <-cmdQ: // 收到控制指令
		if strings.EqualFold(cmd, "quit") {
			log.Println("quit")
			break
		}
	}
}

// 管理所有客户端连接信息； 发送数据
func Dispatcher() {
	connMap := make(map[string]*ConnInfo)
	for {
		select {
		case cmd := <-cmdQ:
			if strings.EqualFold(cmd, "quit") {
				break
			}

		case sendQEle := <-sendQ: // 从发送队列取出 chid->bytes ，找出接收者，然后发送出去。
			senderInfo, found := connMap[sendQEle.chid]
			if !found {
				continue
			}
			if len(senderInfo.recvers) > 0 {
				for _, recvChid := range senderInfo.recvers {
					recverInfo, found := connMap[recvChid]
					if !found {
						continue
					}
					recverInfo.conn.Write(sendQEle.data.Bytes())
				}
			} else {
				for _, v := range connMap {
					v.conn.Write(sendQEle.data.Bytes())
				}
			}

		case connInfo := <-createConnQ: // 接收到连接创建通知时，创建连接信息。
			connMap[connInfo.chid] = connInfo
			log.Println("接收到连接创建通知时，创建连接信息 ", connInfo)

		case chid := <-removeConnQ: // 接收到连接关闭通知时，移除连接信息。
			delete(connMap, chid)
			log.Println("接收到连接关闭通知时，移除连接信息 ", chid)
		}
	}
}

// 创建连接; 创建 Handler
func Accepter(listenAddr string, hsPwd string) {

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
		cmdQ <- "quit"
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go Handler(conn, hsPwd)
	}

}

// 连接建立之后的握手过程
func HandShake(conn net.Conn, hsPwd string) (chid string, recvers []string, err error) {
	err = nil
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		log.Println(scanner.Text()) // Println will add back the final '\n'

		var hsMap map[string]interface{}
		err = json.Unmarshal([]byte(scanner.Text()), &hsMap)
		if err != nil {
			return
		}

		if hsMap["req"] == "hs1" {
			if hsMap["chid"] == nil {
				continue
			}

			if len(hsMap["chid"].(string)) > 0 {
				chid = hsMap["chid"].(string)
			} else {
				// chid =uuid // TODO
			}

			hs2 := map[string]string{"rsp": "hs2", "chid": chid}
			enc := json.NewEncoder(conn)
			enc.Encode(hs2)

		} else if hsMap["req"] == "hs3" {
			if hsMap["recvers"] == nil {
				continue
			}

			r := hsMap["recvers"].([]interface{})
			for _, v := range r {
				recvers = append(recvers, v.(string))
			}

			hs4 := map[string]string{"rsp": "hs4", "result": "OK"}
			enc := json.NewEncoder(conn)
			enc.Encode(hs4)
			break
		}
	}
	if err = scanner.Err(); err != nil {
		log.Println("HandShake scanner:", err)
	}
	return
}

// 客户端连接的处理逻辑。 首先握手，握手失败关闭连接；然后读数据，写入 Dispatcher 的处理队列。
func Handler(conn net.Conn, hsPwd string) {
	defer conn.Close()

	chid, recvers, err := HandShake(conn, hsPwd)
	if err != nil {
		return
	}
	createConnQ <- &ConnInfo{conn, chid, 600, recvers}

	for {
		data := make([]byte, 8192)
		n, err := conn.Read(data) // 读出数据，放入 Dispatcher 的发送队列
		if err != nil {           // 关闭时给 Dispatcher 发送通知
			if err == io.EOF {
				log.Println("Connection closed!")
			} else {
				log.Println("Read error: ", err)
			}
			removeConnQ <- chid
			break
		}
		log.Println("Read data:", data[:n])
		sendQ <- SendQEle{chid, bytes.NewBuffer(data[:n])}
	}
}
