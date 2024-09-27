package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Consumer         sarama.Consumer
	subscriptionRepo SubscriptionRepository
	notifier         NotifierSender
}

func NewConsumer(brokers []string, subscriptionRepo SubscriptionRepository, notifier NotifierSender) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{Consumer: consumer, subscriptionRepo: subscriptionRepo, notifier: notifier}, nil
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
		go c.handleMessage(msg)
	}
}

func (c *Consumer) handleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Received raw message: %s", string(msg.Value))

	var notification NotificationMessage
	err := json.Unmarshal(msg.Value, &notification)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Unmarshaled message: %+v", notification)

	// Проверка подписок
	log.Printf("Fetching subscribers for house ID: %d", notification.HouseID)
	subscribers, err := c.subscriptionRepo.GetSubscribersByHouseID(notification.HouseID)
	if err != nil {
		log.Printf("Failed to get subscribers for house %d: %v", notification.HouseID, err)
		return
	}

	if len(subscribers) == 0 {
		log.Printf("No subscribers found for house ID: %d", notification.HouseID)
		return
	}

	log.Printf("Found %d subscribers for house ID: %d", len(subscribers), notification.HouseID)

	for _, user := range subscribers {
		log.Printf("Sending notification to user ID: %d", user.ID)
		err := c.notifier.SendNotification(user, notification)
		if err != nil {
			log.Printf("Failed to send notification to user %d: %v", user.ID, err)
		}
	}
}
