package kafka

import "github.com/gone-io/gone/v2"

// LoadConsumer is for loading Kafka MQ consumer.
func LoadConsumer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConsumer))
}

// LoaderSyncProducer is for loading Kafka MQ SyncProducer.
func LoaderSyncProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideSyncProducer))
}

// LoaderAsyncProducer is for loading Kafka MQ AsyncProducer.
func LoaderAsyncProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideAsyncProducer))
}
