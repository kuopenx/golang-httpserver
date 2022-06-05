package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	bookListHandler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is a book list.")
		w.WriteHeader(http.StatusOK)

		// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
		for k, v := range r.Header {
			w.Header().Set(k, strings.Join(v, "; "))
		}

		// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		os.Setenv("VERSION", "1.0")
		w.Header().Set("VERSION", os.Getenv("VERSION"))

		// 3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
		fmt.Println("客户端 IP：", r.RemoteAddr)
		fmt.Println("HTTP 返回码：", http.StatusOK)
	}

	// 4. 当访问 localhost/healthz 时，应返回 200
	healthzHandler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "200")
	}

	http.HandleFunc("/book/list", bookListHandler)
	http.HandleFunc("/healthz", healthzHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
