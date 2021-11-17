package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Server struct {
	listener net.Listener
	cons     []net.Conn
	msg      chan string
}

func main() {
	server := initServer()
	go server.DoReceive()
	go server.DoReplay()

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	s := <-ch
	log.Fatal(s)
}

func initServer() *Server {
	listen, err := net.Listen("tcp", ":20200")
	if err != nil {
		log.Error(err)
		return nil
	}
	return &Server{
		listener: listen,
		cons:     []net.Conn{},
		msg:      make(chan string),
	}
}

func (s *Server) DoReceive() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			log.Error(err)
		}
		log.Println(c.RemoteAddr(), c.LocalAddr(), "join server")
		s.cons = append(s.cons, c)
		go func(c net.Conn) {
			scan := bufio.NewScanner(c)
			scan.Split(bufio.ScanLines)
			for scan.Scan() {
				str := scan.Text()
				s.msg <- str
				log.Println("get", fmt.Sprint(c.RemoteAddr(), c.LocalAddr())+" "+str)
			}
			defer c.Close()
		}(c)
	}
}

func (s *Server) DoReplay() {
	for str := range s.msg {
		if strings.HasPrefix(str, "To") {
			id := str[3:7]
			index, err := strconv.Atoi(id)
			if err != nil {
				log.Error(err)
			}
			//s.cons[index].SetWriteDeadline(time.Now().Add(2000))
			n, err := s.cons[index].Write([]byte(str + string('\n')))
			if err != nil {
				log.Error(err)
			}
			log.Println("send to ", index, " "+str, "len is ", n)
			continue
		}

		for i := 0; i < len(s.cons); i++ {
			//s.cons[i].SetWriteDeadline(time.Now().Add(2000))
			s.cons[i].Write([]byte(str + string('\n')))
		}
	}
}
