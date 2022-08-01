package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	metricsNamespace = "httpserver"
	metricsHelp      = "time spent."
	metricsName      = "execution_latency_seconds"
)

var functionLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Name:      metricsName,
		Help:      metricsHelp,
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"step"},
)

type ExecutionTimer struct {
	histogram *prometheus.HistogramVec
	start     time.Time
	last      time.Time
}

func (t *ExecutionTimer) observeTotal() {
	(*t.histogram).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func newExecutionTimer(histogram *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histogram: histogram,
		start:     now,
		last:      now,
	}
}

func registerMetrics() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	registerMetrics()

	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		timer := newExecutionTimer(functionLatency)
		defer timer.observeTotal()

		// 为请求增加模拟延时
		delay := rand.Intn(2001)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Println("延时了", delay, "毫秒")

		io.WriteString(w, "hello.")
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

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	log.Fatal(http.ListenAndServe(":80", nil))
}
