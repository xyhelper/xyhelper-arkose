# xyhelper-arkose

自动获取arkose的token，用于自动化测试

## 1. 安装
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
docker compose up -d
```

或者仅复制 `docker-compose.yml` 中的内容

```yaml
version: '3'
services:
  chat.openai.com:
    image: xyhelper/xyhelper-arkose:latest
    restart: always
    ports:
      - 8199:80
    environment:
      - PORT=80
  chrome:
    image: kasmweb/chrome:1.10.0
    ports:
      - "6901:6901"
    environment:
      - VNC_PW=xyhelper
      - URL=http://chat.openai.com
    shm_size: 512m
```
```bash
docker compose up -d
```

## 2. 使用

### 2.1 获取token
```bash
curl "http://localhost:8199/token"
```

### 2.2 获取token池容量
```bash
curl "http://localhost:8199/ping"
```

### 2.3 主动挂机
```bash
curl "http://localhost:8199/?delay=10"
```

## 3. 增加挂机节点
在节点使用以下 `docker-compose.yml` 启动
```yaml 
version: '3'
services:
  chrome:
    image: kasmweb/chrome:1.10.0
    ports:
      - "6901:6901"
    environment:
      - VNC_PW=xyhelper
      - URL=https://chatarkose.xyhelper.cn  # 修改为你的挂机节点
    shm_size: 512m
```
```bash
docker compose up -d
```