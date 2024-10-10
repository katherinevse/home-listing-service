package handler

type Handler struct {
	tokenManager     TokenManager
	userRepo         UserRepository
	houseRepo        HouseRepository
	flatRepo         FlatRepository
	subscriptionRepo SubscriptionRepository
	producer         ProducerManager
	consumer         ConsumerManager
	logger           Logger
}

func New(
	tokenManager TokenManager,
	userRepo UserRepository,
	houseRepo HouseRepository,
	flatRepo FlatRepository,
	subscriptionRepo SubscriptionRepository,
	producer ProducerManager,
	consumer ConsumerManager,
	logger Logger,
) *Handler {
	return &Handler{
		tokenManager:     tokenManager,
		userRepo:         userRepo,
		houseRepo:        houseRepo,
		flatRepo:         flatRepo,
		subscriptionRepo: subscriptionRepo,
		producer:         producer,
		consumer:         consumer,
		logger:           logger,
	}
}
