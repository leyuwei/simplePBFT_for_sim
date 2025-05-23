package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	nodeId     int
	url        string
	keypair    Keypair
	knownNodes []*KnownNode
	request    *RequestMsg
	replyLog   map[int]*ReplyMsg
	mutex      sync.Mutex
	requestID  int
	filename   string
}

func NewClient(filename string) *Client {
	client := &Client{
		ClientNode.nodeID,
		ClientNode.url,
		KeypairMap[ClientNode.nodeID],
		KnownNodes,
		nil,
		make(map[int]*ReplyMsg),
		sync.Mutex{},
		0,
		filename,
	}
	return client
}

func (c *Client) Start() {
	// Create a timer that triggered every five seconds 
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	// Process the timer task in an individual goroutine
	go func() {
		for range ticker.C {
			c.requestID = c.requestID + 1
			file, err := os.OpenFile(c.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()
			timestamp := time.Now().UnixNano()
			_, err = file.WriteString(fmt.Sprintf("Loop: %d, Timestamp: %d\n", c.requestID, timestamp))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Printf("Loop Times:{%d}\n", c.requestID)
			c.sendRequest()
		}
	}()

	ln, err := net.Listen("tcp", c.url)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go c.handleConnection(conn)
	}
}

func (c *Client) handleConnection(conn net.Conn) {
	req, err := ioutil.ReadAll(conn)
	header, payload, _ := SplitMsg(req)
	if err != nil {
		panic(err)
	}
	switch header {
	case hReply:
		c.handleReply(payload)
	}
}

func (c *Client) sendRequest() {
	msg := fmt.Sprintf("%d work to do!", rand.Int())
	req := Request{
		msg,
		hex.EncodeToString(generateDigest(msg)),
	}
	reqmsg := &RequestMsg{
		"solve",
		int(time.Now().Unix()),
		c.nodeId,
		req,
	}
	sig, err := c.signMessage(reqmsg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	logBroadcastMsg(hRequest, reqmsg)
	send(ComposeMsg(hRequest, reqmsg, sig), c.findPrimaryNode().url)
	c.request = reqmsg
}

func (c *Client) handleReply(payload []byte) {
	var replyMsg ReplyMsg
	err := json.Unmarshal(payload, &replyMsg)
	if err != nil {
		fmt.Printf("error happened:%v", err)
		return
	}
	logHandleMsg(hReply, replyMsg, replyMsg.NodeID)
	c.mutex.Lock()
	c.replyLog[replyMsg.NodeID] = &replyMsg
	rlen := len(c.replyLog)
	c.mutex.Unlock()
	if rlen >= c.countNeedReceiveMsgAmount() {
		//fmt.Println("request success!!")
	}
}

func (c *Client) signMessage(msg interface{}) ([]byte, error) {
	sig, err := signMessage(msg, c.keypair.privkey)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (c *Client) findPrimaryNode() *KnownNode {
	nodeId := ViewID % len(c.knownNodes)
	for _, knownNode := range c.knownNodes {
		if knownNode.nodeID == nodeId {
			return knownNode
		}
	}
	return nil
}

func (c *Client) countTolerateFaultNode() int {
	return (len(c.knownNodes) - 1) / 3
}

func (c *Client) countNeedReceiveMsgAmount() int {
	f := c.countTolerateFaultNode()
	return f + 1
}


