package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Consumer         sarama.Consumer
	SubscriptionRepo SubscriptionRepository // репозиторий для работы с подписками
	//Notifier  NotificationSender     // для отправки уведомлений
}

func NewConsumer(brokers []string, subscriptionRepo SubscriptionRepository) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{Consumer: consumer, SubscriptionRepo: subscriptionRepo}, nil
}
func (c *Consumer) Listen(topic string) {
	log.Printf("Starting consumer for topic: %s", topic)

	partitionConsumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Failed to close partition consumer: %v", err)
		}
	}()

	for msg := range partitionConsumer.Messages() {
		go c.handleMessage(msg) // Обработка в горутине
	}
}

func (c *Consumer) handleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Received raw message: %s", string(msg.Value)) // Логируем сырое сообщение

	var notification NotificationMessage
	err := json.Unmarshal(msg.Value, &notification)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Unmarshaled message: %+v", notification) // Логируем сообщение после десериализации

	// Проверка подписок
	log.Printf("Fetching subscribers for house ID: %d", notification.HouseID)
	subscribers, err := c.SubscriptionRepo.GetSubscribersByHouseID(notification.HouseID)
	if err != nil {
		log.Printf("Failed to get subscribers for house %d: %v", notification.HouseID, err)
		return
	}

	if len(subscribers) == 0 {
		log.Printf("No subscribers found for house ID: %d", notification.HouseID)
		return
	}

	log.Printf("Found %d subscribers for house ID: %d", len(subscribers), notification.HouseID)

}
