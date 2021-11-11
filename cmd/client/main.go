package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
)

func main() {
	createClient()
}

func createClient() {
	conn, err := net.Dial("tcp", ":20200")
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		scan := bufio.NewScanner(conn)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			str := scan.Text()
			log.Println(str)
		}
	}()
	getInToWriter(conn)
	defer conn.Close()
}

func getInToWriter(w io.Writer) {
	for {
		scan := bufio.NewScanner(os.Stdin)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			str := scan.Text()
			log.Println("send", str)
			_, err := w.Write([]byte(str + string('\n')))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
