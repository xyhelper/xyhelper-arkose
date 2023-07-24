# xyhelper-arkose

[ENGLISH](README_EN.md)

自动获取arkose的token，用于自动化测试

## 通知
不再提供P项目BYPASS功能,没有原因,请不要问为什么

## 1. 安装
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
./deploy.sh
```

不要仅复制`docker-compose.yml`，因为`docker-compose.yml`中用到了`Caddyfile`中的配置

## 2. 使用

### 2.1 获取token
```bash
curl "http://服务器IP:8199/token"
```

### 2.2 获取token池容量
```bash
curl "http://服务器IP:8199/ping"
```

### 2.3 主动挂机
```bash
curl "http://服务器IP:8199/?delay=10"
```

## 3. 增加挂机节点
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
```

修改`docker-compose.yml` 取消   # - FORWORD_URL=https://chatarkose.xyhelper.cn/pushtoken 的注释

执行`./deploy.sh`

## 4. 管理chrome

登陆地址：https://服务器IP:6901

用户名：kasm_user

默认密码：xyhelper  

## 5. 公共节点

获取token地址：https://chatarkose.xyhelper.cn/token

查询token池容量：https://chatarkose.xyhelper.cn/ping