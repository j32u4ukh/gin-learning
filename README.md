# RealMoneySo

## Introduction

* DB は MySQL
* アプリ用サーバーは AL2023
* Go は最新バージョン

## Usage

### docker-compose.ymlの更新を反映させる

```bash
docker-compose up -d
```

### Dockerfileの更新を反映させる

```bash
docker-compose up -d --build
```

### リビルド

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### app コンテナに入る

```bash
docker-compose exec app bash
```
