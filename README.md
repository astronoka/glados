# GLaDOS
Github Lifeform and Disk Operating System

## memo

### build glados-server container

```bash
make
```

### start glados-server in container (auto rebuild)

```bash
DOCKER_NETWORK=some_docker_network make run-glados-server
```

### start glados-server (auto rebuild)

```bash
reflex -r '\.go$' -s -- sh -c 'go run -v example/cmd/glados-server/main.go'
```

## .env

| key | description |
| --- | --- |
| PORT | listen port number |
| BOT_NAME | bot name |
| GLADOS_GITHUB_NOTIFIER_SECRET | github webhook secret string |
| GLADOS_SLACK_BOT_UAER_TOKEN | slack bot user token |
| GLADOS_DATASTORE_MYSQL_DSN | mysql storage dsn (user:password@tcp(127.0.0.1:3306)/glados?parseTime=true) |
