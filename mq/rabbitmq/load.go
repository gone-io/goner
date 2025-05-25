package rabbitmq

import "github.com/gone-io/gone/v2"

// LoadConnection is for loading RabbitMQ connection.
func LoadConnection(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConnection))
}

// LoadChannel is for loading RabbitMQ channel.
func LoadChannel(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideChannel))
}

// LoadProducer is for loading RabbitMQ producer.
func LoadProducer(loader gone.Loader) error {
	return loader.
		MustLoadX(LoadConnection).
		MustLoadX(LoadChannel).
		Load(gone.WrapFunctionProvider(ProvideProducer), gone.IsDefault(new(IProducer)))
}

// LoadConsumer is for loading RabbitMQ consumer.
func LoadConsumer(loader gone.Loader) error {
	return loader.
		MustLoadX(LoadConnection).
		MustLoadX(LoadChannel).
		Load(gone.WrapFunctionProvider(ProvideConsumer), gone.IsDefault(new(IConsumer)))
}

// LoadAll is for loading all RabbitMQ components.
func LoadAll(loader gone.Loader) error {
	loader.
		MustLoadX(LoadProducer).
		MustLoadX(LoadConsumer)
	return nil
}
