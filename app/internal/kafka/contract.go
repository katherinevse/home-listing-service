package kafka

import "app/internal/repository/model"

type SubscriptionRepository interface {
	CreateSubscription(subscriber *model.Subscription) error
	GetSubscribersByHouseID(houseID int) ([]model.User, error)
}

type NotifierSender interface {
	SendNotification(user model.User, notification NotificationMessage) error
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}
