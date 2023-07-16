# xyhelper-arkose

自动获取arkose的token，用于自动化测试

## 1. 安装
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
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

