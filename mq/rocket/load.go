package rocket

import "github.com/gone-io/gone/v2"

// LoadConsumer is for loading Rocket MQ consumer.
func LoadConsumer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConsumer))
}

// LoaderProducer is for loading Rocket MQ producer.
func LoaderProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideProducer))
}
