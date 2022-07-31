package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	booklistHandler := func(w http.ResponseWriter, r *http.Request) {
		// 为请求增加模拟延时
		delay := rand.Intn(3)
		time.Sleep(time.Millisecond * time.Duration(delay))
		fmt.Println("延时了", delay, "秒")

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

	// 处理metrics
	metricsHandler := func(w http.ResponseWriter, r *http.Request) {
		handler := promhttp.Handler()
		handler.ServeHTTP(w, r)
	}
	http.HandleFunc("/booklist", booklistHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/metrics", metricsHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
