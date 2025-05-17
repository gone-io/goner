package kafka

import "github.com/gone-io/gone/v2"

// LoadConsumer is for loading sarama.Consumer.
// Consumer manages PartitionConsumers which process Kafka messages from brokers. You MUST call Close()
// on a consumer to avoid leaks, it will not be garbage-collected automatically when it passes out of
// scope.
func LoadConsumer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConsumer))
}

// LoadConsumerGroup is for loading sarama.ConsumerGroup.
// ConsumerGroup is responsible for dividing up processing of topics and partitions
// over a collection of processes (the members of the consumer group).
func LoadConsumerGroup(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideConsumerGroup))
}

// LoaderSyncProducer is for loading sarama.SyncProducer.
func LoaderSyncProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideSyncProducer))
}

// LoaderAsyncProducer is for loading sarama.AsyncProducer.
func LoaderAsyncProducer(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(ProvideAsyncProducer))
}
