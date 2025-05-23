package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

var urlName = "localhost:%d"

type Server struct {
	node     *Node
	url      string
	filename string
}

func nodeIdToPort(nodeId int) int {
	return nodeId + 8180
}

func NewServer(nodeId int, filename string) *Server {
	filename = strconv.Itoa(nodeId) + filename
	server := &Server{
		NewNode(nodeId, filename),
		fmt.Sprintf(urlName, nodeIdToPort(nodeId)),
		filename,
	}
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}
	return server
}

func (s *Server) Start() {
	s.node.Start()
	ln, err := net.Listen("tcp", s.url)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("server start at %s\n", s.url)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	req, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	s.node.msgQueue <- req
}
