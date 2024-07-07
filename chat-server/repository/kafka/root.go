package kafka

import (
	"chat-server/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	config   *config.Config
	producer *kafka.Producer
}

func NewKafka(config *config.Config) (*Kafka, error) {
	k := &Kafka{config: config}
	var err error

	if k.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.URL,
		"client.id":         config.Kafka.ClientID,
		"acks":              "all", // 메시지 전송 시 고가용성을 위해 복제본을 어디까지 저장할지에 대한 설정 값
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func (k *Kafka) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, ch); err != nil {
		return nil, err
	} else {
		return <-ch, nil
	}
}
