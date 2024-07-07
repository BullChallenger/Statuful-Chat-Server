package kafka

import (
	"chat-controller/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	config   *config.Config
	consumer *kafka.Consumer
}

func NewKafka(config *config.Config) (*Kafka, error) {
	k := &Kafka{config: config}
	var err error

	if k.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.URL,
		"group.id":          config.Kafka.GroupID,
		"auto.offset.reset": "latest",
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func (k *Kafka) Poll(timeoutMs int) kafka.Event {
	return k.consumer.Poll(timeoutMs)
}

func (k *Kafka) RegisterSubTopic(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		return err
	} else {
		return nil
	}
}
