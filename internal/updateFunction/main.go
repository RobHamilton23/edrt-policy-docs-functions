package update

import (
	"context"
	"encoding/json"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

type UpdateHandler struct {
	logger *logrus.Logger
	store  *firestore.Firestore
}

func NewUpdateHandler(logger *logrus.Logger, fs *firestore.Firestore) UpdateHandler {
	return UpdateHandler{
		logger: logger,
		store:  fs,
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
		logger.Errorf("event.DataAs: %w", err)

		// Bad messages should not trigger the cloud function to retry
		return nil
	}

	msgText := string(msg.Message.Data)

	var pdocsMessage types.PolicyDocsMessage
	err := json.Unmarshal([]byte(msgText), &pdocsMessage)
	if err != nil {
		logger.Errorf(
			"unable to parse pubsub message (%s) as PolicyDocsMessage: %w",
			msgText,
			err,
		)

		// When the message fails to deserialize, we don't want the funciton to
		// repeatedly retry, so we'll just exit as if everything's fine and the
		// bad message will be ignored.
		return nil
	}

	logger.WithField("PubSub", true).WithField("Data", msgText).Info("Pubsub triggered")
	logger.Info("Iterating over attributes")
	for x := range msg.Message.Attributes {
		logger.WithField(x, msg.Message.Attributes[x]).Info("Attribute")
	}
	logger.Info("Iteration complete")

	result := u.store.Read(ctx)
	logger.WithField("Data", result).Info("Data from firestore")
	return nil
}
