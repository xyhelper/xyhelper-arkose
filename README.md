# xyhelper-arkose

---
# 因OPENAI修改了生成的token的有效期，token池模式已无存在意义，本项目终结.
---
[ENGLISH](README_EN.md)

自动获取 arkose 的 token，用于自动化测试

## 通知

不再提供 P 项目 BYPASS 功能,没有原因,请不要问为什么

## 1. 安装

创建`docker-compose.yml`文件

```yaml
version: "3"
services:
  broswer:
    image: xyhelper/xyhelper-arkose:latest
    ports:
      - "8199:3000"
```

启动

```bash
docker-compose up -d
```

## 2. 使用

### 2.1 获取 token

```bash
curl "http://服务器IP:8199/token"
```

### 2.2 获取 token 池容量(废弃)

```bash
curl "http://服务器IP:8199/ping"
```

### 2.3 查看当前 payload

```bash
curl "http://服务器IP:8199/payload"
```

## 3. 公共节点

获取 token 地址：https://chatarkose.xyhelper.cn/token

## 加入交流群

![微信](https://xyhelper.github.io/xyhelperkf.png)
