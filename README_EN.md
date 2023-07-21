# xyhelper-arkose

Automatically obtain Arkose tokens for automated testing.

## 1. Installation
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
docker compose up -d
```

Alternatively, copy only the contents of `docker-compose.yml`

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

## 2. Usage

### 2.1 Get Token
```bash
curl "http://localhost:8199/token"
```

### 2.2 Get Token Pool Capacity
```bash
curl "http://localhost:8199/ping"
```

### 2.3 Hang Up Actively
```bash
curl "http://localhost:8199/?delay=10"
```

## 3. Adding Hanging Nodes
Start the node with the following `docker-compose.yml`
```yaml 
version: '3'
services:
  chrome:
    image: kasmweb/chrome:1.10.0
    ports:
      - "6901:6901"
    environment:
      - VNC_PW=xyhelper
      - URL=https://chatarkose.xyhelper.cn  # Change this to your hanging node
    shm_size: 512m
```
```bash
docker compose up -d
```

## 4. Managing Chrome

Login URL: https://localhost:6901

Username: kasm_user

Default Password: xyhelper