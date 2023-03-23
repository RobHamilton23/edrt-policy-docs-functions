package update

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

type UpdateHandler struct {
	logger *logrus.Logger
}

func NewUpdateHandler(logger *logrus.Logger) UpdateHandler {
	return UpdateHandler{
		logger: logger,
	}
}

func (u *UpdateHandler) PolicyDocUpdated(ctx context.Context, e event.Event) error {
	/**
	* This function is intended to be the pub/sub interface. It will do as little as
	* possible. It will call into other packages to read from firestore, denormalize
	* the data, then write to firestore
	**/

	logger := u.logger.WithField("Function", "PolicyDocUpdated")

	var msg types.MessagePublishedData

	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	msgText := string(msg.Message.Data)

	logger.WithField("PubSub", true).WithField("Data", msgText).Info("Pubsub triggered")
	logger.Info("Iterating over attributes")
	for x := range msg.Message.Attributes {
		logger.WithField(x, msg.Message.Attributes[x]).Info("Attribute")
	}
	logger.Info("Iteration complete")
	return nil
}
