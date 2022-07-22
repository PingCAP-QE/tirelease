package gconsumer

import "github.com/google/go-github/v41/github"

type PREventConsumer interface {
	Validate(event github.PullRequestEvent) bool
	Consume(event github.PullRequestEvent) error
}

var pREventConsumer []PREventConsumer

func GetPREventConsumers() []PREventConsumer {
	if pREventConsumer == nil {
		initPREventConsumerMap()
	}

	return pREventConsumer
}

func initPREventConsumerMap() {
	pREventConsumer = make([]PREventConsumer, 0)
}
