package kafka

import "app/internal/repository/model"

type SubscriptionRepository interface {
	CreateSubscription(subscriber *model.Subscription) error
	GetSubscribersByHouseID(houseID int) ([]model.User, error)
}

type NotifierSender interface {
	SendNotification(user model.User, notification NotificationMessage) error
}
