package rocket

//go:generate mockgen -destination=rocket_mock.go -package=rocket github.com/apache/rocketmq-clients/golang/v5 SimpleConsumer,Producer
