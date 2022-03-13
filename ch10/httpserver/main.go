package main

import (
	"cncamphomework/ch10/httpserver/metrics"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wonderivan/logger"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

/*
1.接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 localhost/healthz 时，应返回200
*/

func main() {
	fmt.Println("http server start.")
	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/", roothandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/delay", delay)
	mux.Handle("/metrics",promhttp.Handler())
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		//fmt.Println("Http server error")
		fmt.Println(err)
		return
	}
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	rsp, err := http.Get(fmt.Sprintf("http://[%s]:80", r.Host))
	if err != nil {
		//fmt.Println("Http server error")
		fmt.Println(err)
		return
	}
	io.WriteString(w, "VERSION:"+string(os.Getenv("VERSION"))+"\n")                            // 返回系统环境变量VERSION
	io.WriteString(w, fmt.Sprintf("Client IP: %s \nStatusCode: %d\n", r.Host, rsp.StatusCode)) // 输出客户端IP, http返回码
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "StatusCode:", http.StatusOK) // healthz handle 返回200
}

func delay(w http.ResponseWriter, r *http.Request) {
		timer := metrics.NewTimer()
		defer timer.ObserverTotal()
		delay := randInt(10, 2000)
		time.Sleep(time.Millisecond * time.Duration(delay))
		roothandler(w, r)
		io.WriteString(w, "Wait for Delay")
		logger.Info("Respond in %d ms", delay)
	}


func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}