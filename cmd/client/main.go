package main

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {

}

func tailServer(values ...string) {
	conn, err := net.Dial("tcp", ":20200")
	if err != nil {
		log.Error(err)
		createServer()
		return
	}
	if len(values) == 0 {
		values = append(values, "default")
	}
	for i := 0; i < len(values); i++ {
		_, err = conn.Write([]byte(values[i] + string('\n')))
	}
	if err != nil {
		log.Error(err)
	}
	conn.Close()
}

func createServer() {
	listen, err := net.Listen("tcp", ":20200")
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		for {
			c, err := listen.Accept()
			if err != nil {
				log.Error(err)
			}
			log.Println(c.RemoteAddr(), c.LocalAddr())
			scan := bufio.NewScanner(c)
			scan.Split(bufio.ScanLines)
			for scan.Scan() {
				str := scan.Text()
				log.Println(str)
			}
			c.Close()
		}
	}()
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, signals...)
	s := <-ch
	log.Fatal(s)
}
