package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

type Consumer struct {
	Consumer         sarama.Consumer
	subscriptionRepo SubscriptionRepository
	notifier         NotifierSender
	logger           Logger
}

func NewConsumer(brokers []string, subscriptionRepo SubscriptionRepository, notifier NotifierSender, logger *slog.Logger) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{Consumer: consumer, subscriptionRepo: subscriptionRepo, notifier: notifier, logger: logger}, nil
}
func (c *Consumer) Listen(topic string) {
	c.logger.Info("Starting consumer for topic", "topic", topic)

	partitionConsumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		c.logger.Error("Failed to start consumer", "error", err) //TODO фатал
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			c.logger.Warn("Failed to close partition consumer", "error", err)
		}
	}()

	for msg := range partitionConsumer.Messages() {
		go c.handleMessage(msg)
	}
}

func (c *Consumer) handleMessage(msg *sarama.ConsumerMessage) {
	c.logger.Debug("Received raw message", "message", string(msg.Value))

	var notification NotificationMessage
	err := json.Unmarshal(msg.Value, &notification)
	if err != nil {
		c.logger.Warn("Failed to unmarshal message", "error", err)
		return
	}

	c.logger.Debug("Unmarshaled message", "notification", notification)

	// Проверка подписок
	c.logger.Info("Fetching subscribers for house ID", "houseID", notification.HouseID)
	subscribers, err := c.subscriptionRepo.GetSubscribersByHouseID(notification.HouseID)
	if err != nil {
		c.logger.Warn("Failed to get subscribers for house", "houseID", notification.HouseID, "error", err)
		return
	}

	if len(subscribers) == 0 {
		c.logger.Info("No subscribers found for house ID", "houseID", notification.HouseID)
		return
	}
	c.logger.Info("Found subscribers", "count", len(subscribers), "houseID", notification.HouseID)

	for _, user := range subscribers {
		c.logger.Info("Sending notification to user", "userID", user.ID)
		err := c.notifier.SendNotification(user, notification)
		if err != nil {
			c.logger.Warn("Failed to send notification", "userID", user.ID, "error", err)
		}
	}
}
