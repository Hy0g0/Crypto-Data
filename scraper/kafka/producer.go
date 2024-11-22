package kafka

import (
    "encoding/json"
    "github.com/Shopify/sarama"
    "os"
)

type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &Producer{
        producer: producer,
        topic:    topic,
    }, nil
}

func (p *Producer) SendMessage(data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: p.topic,
        Value: sarama.StringEncoder(jsonData),
    }

    _, _, err = p.producer.SendMessage(msg)
    return err
} 