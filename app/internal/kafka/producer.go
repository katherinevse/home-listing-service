package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: producer}, nil
}

// PublishNotification - отправляет уведомление о новой квартире в топик
func (p *Producer) PublishNotification(houseID int, flatNumber int, message string) error {
	notification := NotificationMessage{
		HouseID:    houseID,
		FlatNumber: flatNumber,
		Message:    message,
	}
	value, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "new-flat",
		Value: sarama.StringEncoder(value),
	}

	_, _, err = p.Producer.SendMessage(msg)
	return err
}
