模块三作业：

---

- 构建本地镜像
- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
- 将镜像推送至 docker 官方镜像仓库
- 通过 docker 命令本地启动 httpserver
- 通过 nsenter 进入容器查看 IP 配置

---

1. 编写dockerfile文件

   ```dockerfile
   FROM golang:1.17.6-alpine AS build
   
   WORKDIR /go/src/httpserver/
   
   COPY httpserver/* /go/src/httpserver/
   RUN go env -w GO111MODULE=auto && go build -o /bin/httpserver
   
   FROM alpine:3.15.0
   COPY --from=build /bin/httpserver /bin/
   EXPOSE 80
   ENTRYPOINT ["/bin/httpserver"]
   ```

2. 构建本地镜像

   > docker build -t fingerf/httpserver:1.0 .

3. 通过 docker 命令本地启动 httpserver

   > sudo docker run -it -p 80:80 --name="httpserver" fingerf/httpserver:1.0

4. 通过 nsenter 进入容器查看 IP 配置

   ```sh
   ➜  ~ sudo docker inspect httpserver | grep -i pid                               
               "Pid": 5886,                                                        
               "PidMode": "",                                                      
               "PidsLimit": null,            
   
   ➜  ~ sudo ls -la /proc/5886/ns/
   total 0
   dr-x--x--x 2 root root 0 Jan 16 08:06 .
   dr-xr-xr-x 9 root root 0 Jan 16 08:06 .. 
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 cgroup -> 'cgroup:[4026531835]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 ipc -> 'ipc:[4026532641]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 mnt -> 'mnt:[4026532639]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:06 net -> 'net:[4026532644]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 pid -> 'pid:[4026532642]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 pid_for_children -> 'pid:[4026532642]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 user -> 'user:[4026531837]'
   lrwxrwxrwx 1 root root 0 Jan 16 08:11 uts -> 'uts:[4026532640]'
   
   ➜  ~ sudo nsenter -t 5886 -n ip add
   1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
       link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
       inet 127.0.0.1/8 scope host lo
          valid_lft forever preferred_lft forever
   18: eth0@if19: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
       link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
       inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
          valid_lft forever preferred_lft forever
   ```

5. 将镜像推送至 docker 官方镜像仓库

   ``` sh
   ➜  ~ sudo docker login
   Login Succeeded
   ➜  ~ sudo docker push fingerf/httpserver:1.0
   The push refers to repository [docker.io/fingerf/httpserver]
   a7cdea4ec15b: Pushed 
   8d3ac3489996: Mounted from library/golang 
   1.0: digest: sha256:e7ae2e82677958776878da8e126b609725add0dd4b3a3b8256695140c5a167b2 size: 739
   ```

   