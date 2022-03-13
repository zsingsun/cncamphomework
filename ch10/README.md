模块三作业：

---

- 为 HTTPServer 添加 0-2 秒的随机延时；
- 为 HTTPServer 项目添加延时 Metric；
- 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；
- 从 Promethus 界面中查询延时指标数据；
- （可选）创建一个 Grafana Dashboard 展现延时分配情况。

---

1. 修改Deployment, image: fingerf/http-server:1.1-metrics

2. 新增 http-server-prometheus-deployment.yaml

3. 安装 helm

4. 安装 loki-stack

   ```shell
   ### add grafana repo
   helm repo add grafana https://grafana.github.io/helm-charts
   ```

   ```shell
   ### install loki-stack
   helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
   ```

   ```shell
   ### if you get the following error, that means your k8s version is too new to install
   Error: unable to build kubernetes objects from release manifest: [unable to recognize "": no matches for kind "ClusterRole" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "ClusterRoleBinding" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "Role" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "RoleBinding" in version "rbac.authorization.k8s.io/v1beta1"]
   ```

   ```shell
   ### download loki-stack
   helm pull grafana/loki-stack
   tar -xvf loki-stack-2.5.0.tgz
   ```

   ```shell
   cd loki-stack
   
   ### replace all `rbac.authorization.k8s.io/v1beta1` with `rbac.authorization.k8s.io/v1` by
   ### 需替换该目录及其子目录下的所有文件中的指定字符串
   sed -i s#rbac.authorization.k8s.io/v1beta1#rbac.authorization.k8s.io/v1#g `grep rbac.authorization.k8s.io/v1beta1 -rl ./`
   ```

   ```
   ### install loki locally
   helm upgrade --install loki ./loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
   ```

5. 更新 loki-grafana service 和 loki-prometheus-server service 为 NodePort type

6. 查看 grafana 的登陆信息 

7. 部署 http-server prometheus deployment

   ```shell
   k apply -f yaml/http-server-prometheus-deployment.yaml
   ```

8. 访问指标生成的接口

   ```shell
   curl 192.168.11.230:9000/api/apps/v2/index
   ```

9. 访问指标采集的接口

   ```shell
   curl 192.168.11.230:9000/metrics
   ```

10. promethemus 查看指标采集数据

11. import httpserver-latency.json 配置，查看指标采集数据