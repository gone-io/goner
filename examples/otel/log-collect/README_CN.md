```bash
docker compose up -d

go mod init examples/otel/collect

gonectl install goner/otel/log/http
gonectl install goner/zap
gonectl install goner/viper
```