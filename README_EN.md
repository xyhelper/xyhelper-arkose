# xyhelper-arkose

[ENGLISH](README_EN.md)

Automatically acquire tokens for arkose for automated testing purposes.

## Notification

The BYPASS feature for Project P is no longer available, without any specific reasons provided. Please refrain from inquiring about the reason.

## 1. Installation

Create a `docker-compose.yml` file

```yaml
version: "3"
services:
  broswer:
    image: xyhelper/xyhelper-arkose:latest
    ports:
      - "8199:3000"
```

Start

```bash
docker-compose up -d
```

## 2. Usage

### 2.1 Obtain token

```bash
curl "http://serverIP:8199/token"
```

### 2.2 Obtain token pool capacity (Deprecated)

```bash
curl "http://serverIP:8199/ping"
```

### 2.3 View current payload

```bash
curl "http://serverIP:8199/payload"
```

## 3. Public Nodes

Token acquisition address: https://chatarkose.xyhelper.cn/token

Token pool capacity inquiry: https://chatarkose.xyhelper.cn/ping