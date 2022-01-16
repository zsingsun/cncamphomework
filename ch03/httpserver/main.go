package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
1.接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 localhost/healthz 时，应返回200
*/

func main() {
	fmt.Println("http server start.")
	http.HandleFunc("/", roothandler)
	http.HandleFunc("/healthz", healthz)
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
