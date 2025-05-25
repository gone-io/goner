package rocket

import "github.com/gone-io/gone/v2"

// LoadConsumer is for loading Rocket MQ consumer.
func LoadConsumer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConsumer))
}

// LoadProducer is for loading Rocket MQ producer.
func LoadProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideProducer))
}
