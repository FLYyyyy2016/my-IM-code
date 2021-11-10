package main

import (
	"testing"
)

func TestServer(t *testing.T) {
	createServer()
}

func TestClient(t *testing.T) {
	tailServer("baidu", "ceshi", "你好")
}
