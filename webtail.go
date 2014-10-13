package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var userFile string = ""

func logHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "max-age=0")
	//userFile := "test.txt"
	fin, err := os.Open(userFile)
	defer fin.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stat, _ := fin.Stat() //获取文件状态
	old_filesize := stat.Size()
	fin.Seek(old_filesize-100, 0)
	br := bufio.NewReader(fin)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		_, werr := w.Write([]byte(fmt.Sprintf("%s%s", line, "<br>")))
		if werr != nil {
			fmt.Println("http over")
			break
		}
		if f, ok := w.(http.Flusher); ok {
			fmt.Println("flush")
			f.Flush()
		}
	}
}

func main() {
	port := flag.String("p", "8888", "port")
	log_file := flag.String("f", "test.txt", "log file")
	flag.Parse()
	userFile = *log_file
	fmt.Printf("port: %s\n", *port)
	fmt.Printf("log file : %s\n", userFile)

	http.HandleFunc("/", logHandler)
	http.ListenAndServe(":"+*port, nil)
}
