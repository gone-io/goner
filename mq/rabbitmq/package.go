package rabbitmq

//go:generate mockgen -destination=rabbitmq_mock.go -package=rabbitmq github.com/rabbitmq/amqp091-go Connection,Channel
