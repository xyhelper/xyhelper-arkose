# xyhelper-arkose


Automatically fetches tokens for arkose to enable automated testing.

## 1. Installation
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
./deploy.sh
```

Do not only copy `docker-compose.yml`, as it relies on configurations from `Caddyfile`.

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
```bash
git clone https://github.com/xyhelper/xyhelper-arkose.git
cd xyhelper-arkose
```

Modify `docker-compose.yml` and remove the '#' from the line: `- FORWORD_URL=https://chatarkose.xyhelper.cn/pushtoken`.

Execute `./deploy.sh`

## 4. Managing Chrome

Login URL: https://localhost:6901

Username: kasm_user

Default Password: xyhelper

## 5. Public Nodes

Get Token URL: https://chatarkose.xyhelper.cn/token

Check Token Pool Capacity: https://chatarkose.xyhelper.cn/ping