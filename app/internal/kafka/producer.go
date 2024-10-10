package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.SyncProducer
	logger   Logger
}

func NewProducer(brokers []string, logger Logger) (*Producer, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: producer, logger: logger}, nil
}

// PublishNotification отправляет уведомление о новой квартире в топик
func (p *Producer) PublishNotification(houseID int, flatNumber int, message string) error {
	notification := NotificationMessage{
		HouseID:    houseID,
		FlatNumber: flatNumber,
		Message:    message,
	}
	value, err := json.Marshal(notification)
	if err != nil {
		p.logger.Error("Failed to marshal notification", "error", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "new-flat",
		Value: sarama.StringEncoder(value),
	}

	// Отправляем сообщение и обрабатываем ошибки
	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send message to Kafka", "error", err)
		return err
	}

	p.logger.Info("Message sent to Kafka", "topic", msg.Topic, "partition", partition, "offset", offset)
	return nil
}
